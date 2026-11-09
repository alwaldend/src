package archive

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/archivepb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	"google.golang.org/protobuf/proto"
)

const (
	// MaxBytes bounds an archive object before ZIP directory decoding.
	MaxBytes int64 = 1 << 30
	// maxEntries bounds ZIP directory allocation, including ignored media entries.
	maxEntries = 10000
	// maxEntryBytes bounds one selected data file before JSON decoding.
	maxEntryBytes = 64 << 20
	// maxDecodedBytes bounds selected uncompressed data across multipart exports.
	maxDecodedBytes = 256 << 20
	// maxRecordBytes bounds one tweet's encoded data.
	maxRecordBytes = 2 << 20
	// maxRecords bounds normalized activity records.
	maxRecords = 50000
	// maxInventoryBytes bounds retained activity text and metadata input.
	maxInventoryBytes = 32 << 20
)

// contextReader interrupts decompression between bounded underlying reads.
type contextReader struct {
	// ctx belongs to the evaluation job.
	ctx context.Context
	// source is one selected ZIP entry.
	source io.Reader
}

// Read stops reading an entry after evaluation cancellation.
func (r contextReader) Read(data []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, fmt.Errorf("read archive data: %w", err)
	}
	n, err := r.source.Read(data)
	if err != nil && err != io.EOF {
		return n, fmt.Errorf("read ZIP entry: %w", err)
	}
	return n, err
}

// archiveReader exposes original file data and one immutable validated directory.
// Its canonical footer excludes comments so zip.NewReader cannot select another EOCD.
type archiveReader struct {
	// ctx interrupts metadata and entry reads during standard-library parsing.
	ctx context.Context
	// source retains random access to the original local-file data.
	source io.ReaderAt
	// offset starts the immutable directory and canonical footer.
	offset int64
	// metadata contains at most 8 MiB of headers plus the 22-byte footer.
	metadata []byte
}

// ReadAt checks cancellation and prevents metadata reads from returning unchecked bytes.
func (r *archiveReader) ReadAt(data []byte, position int64) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, fmt.Errorf("read archive object: %w", err)
	}
	if position < 0 {
		return 0, fmt.Errorf("negative archive offset")
	}
	total := 0
	if position < r.offset {
		length := int(min(int64(len(data)), r.offset-position))
		n, err := r.source.ReadAt(data[:length], position)
		total, position = n, position+int64(n)
		if err != nil && err != io.EOF {
			return n, fmt.Errorf("read archive file data: %w", err)
		}
		if n < length {
			return n, io.EOF
		}
	}
	if total < len(data) {
		start := position - r.offset
		if start >= int64(len(r.metadata)) {
			return total, io.EOF
		}
		total += copy(data[total:], r.metadata[start:])
	}
	if total < len(data) {
		return total, io.EOF
	}
	return total, nil
}

// Export provides sequential, bounded reads of literal YTD datasets from one ZIP.
// Its caller owns the underlying source and must keep it open through evaluation.
type Export struct {
	// AccountID is the identity validated from data/account.js.
	AccountID string
	// account retains the validated typed identity and optional profile fields.
	account *archivepb.Account
	// files indexes validated ZIP entries without extracting them.
	files map[string]*zip.File
	// decoded counts actual decompressed bytes across all requested files.
	decoded int64
	// retained bounds cached JSON rows across every requested dataset.
	retained int
	// records bounds the aggregate number of cached rows.
	records int
	// cache contains immutable JSON rows; Lua receives fresh projections.
	cache map[string][]json.RawMessage
	// failure records any unsuccessful read even if Lua catches the immediate error.
	failure error
}

// Datasets lists dataset names without decoding their contents.
func (e *Export) Datasets() []string {
	names := []string{}
	for name, file := range e.files {
		if file.Mode().IsRegular() && strings.HasPrefix(name, "data/") && strings.HasSuffix(name, ".js") && !strings.Contains(strings.TrimPrefix(name, "data/"), "/") {
			category := strings.TrimSuffix(strings.TrimPrefix(name, "data/"), ".js")
			if !strings.Contains(category, "-part") {
				names = append(names, category)
			}
		}
	}
	sort.Strings(names)
	return names
}

// Failure retains a dataset-read failure across protected Lua calls.
func (e *Export) Failure() error {
	if e.failure != nil {
		return fmt.Errorf("archive dataset read failed: %w", e.failure)
	}
	return nil
}

// Account returns the validated account without the YTD envelope.
func (e *Export) Account(ctx context.Context) (*archivepb.Account, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("read typed account: %w", err)
	}
	return proto.Clone(e.account).(*archivepb.Account), nil
}

// Following returns validated followed-account records for the Lua dataset API.
func (e *Export) Following(ctx context.Context) ([]*archivepb.Following, error) {
	rows, err := typedRows(ctx, e, "following", parseFollowing)
	if err != nil {
		return nil, fmt.Errorf("load followed accounts: %w", err)
	}
	return rows, nil
}

// Likes returns validated liked-post records for the Lua dataset API.
func (e *Export) Likes(ctx context.Context) ([]*archivepb.Like, error) {
	rows, err := typedRows(ctx, e, "like", parseLike)
	if err != nil {
		return nil, fmt.Errorf("load liked posts: %w", err)
	}
	return rows, nil
}

// Tweets imports and validates protobuf records without another archive representation.
func (e *Export) Tweets(ctx context.Context) (result []*archivepb.Tweet, resultErr error) {
	defer func() {
		if resultErr != nil && e.failure == nil {
			e.failure = resultErr
		}
	}()
	records, err := typedRows(ctx, e, "tweets", func(row json.RawMessage) (*archivepb.Tweet, error) {
		value, err := parseTweet(row, e.AccountID)
		if err != nil {
			return nil, fmt.Errorf("parse account tweet: %w", err)
		}
		return value, nil
	})
	if err != nil {
		return nil, fmt.Errorf("import tweets: %w", err)
	}
	seen, versions := map[string]bool{}, map[string]string{}
	for _, value := range records {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("validate tweet inventory: %w", err)
		}
		if seen[value.Id] {
			return nil, fmt.Errorf("duplicate export activity")
		}
		seen[value.Id] = true
		for _, version := range append([]string{value.Id}, value.EditHistoryIds...) {
			if !model.ValidID(version) || (versions[version] != "" && versions[version] != value.Id) {
				return nil, fmt.Errorf("invalid or conflicting edit history")
			}
			versions[version] = value.Id
		}
	}
	return records, nil
}

// Read returns cached immutable raw rows from consecutive dataset parts.
// Callers must not mutate returned bytes; interpreter projections are independent.
func (e *Export) Read(ctx context.Context, category string) (rows []json.RawMessage, resultErr error) {
	defer func() {
		if resultErr != nil && e.failure == nil {
			e.failure = resultErr
		}
	}()
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("before dataset read: %w", err)
	}
	if category == "" || len(category) > 128 || strings.IndexFunc(category, func(c rune) bool {
		return !(c >= 'a' && c <= 'z') && !(c >= 'A' && c <= 'Z') && !(c >= '0' && c <= '9') && c != '_' && c != '-'
	}) >= 0 || strings.Contains(category, "-part") {
		return nil, fmt.Errorf("invalid archive dataset name")
	}
	if cached, ok := e.cache[category]; ok {
		return cached, nil
	}
	prefix := "data/" + category
	first := e.files[prefix+".js"]
	if first == nil || !first.Mode().IsRegular() {
		return nil, fmt.Errorf("archive dataset not found")
	}
	parts := map[int]*zip.File{0: first}
	for name, file := range e.files {
		if strings.HasPrefix(name, prefix+"-part") && strings.HasSuffix(name, ".js") {
			number := strings.TrimSuffix(strings.TrimPrefix(name, prefix+"-part"), ".js")
			index, err := strconv.Atoi(number)
			if err != nil || index < 1 || index > maxEntries || strconv.Itoa(index) != number || !file.Mode().IsRegular() {
				return nil, fmt.Errorf("invalid dataset part")
			}
			parts[index] = file
		}
	}
	rows = []json.RawMessage{}
	for index := 0; index < len(parts); index++ {
		file := parts[index]
		if file == nil {
			return nil, fmt.Errorf("missing dataset part")
		}
		part, err := readPart(ctx, file, strings.ReplaceAll(category, "-", "_"), index, &e.decoded)
		if err != nil {
			return nil, fmt.Errorf("decode dataset part: %w", err)
		}
		for _, row := range part {
			if e.retained > maxInventoryBytes-len(row) || e.records >= maxRecords {
				return nil, fmt.Errorf("archive dataset bounds exceeded")
			}
			e.retained += len(row)
			e.records++
			rows = append(rows, row)
		}
	}
	e.cache[category] = rows
	return rows, nil
}

// Open validates ZIP metadata and account identity without reading other datasets.
// The caller must retain source until all dataset reads have completed.
func Open(ctx context.Context, source io.ReaderAt, size int64) (*Export, error) {
	if size < 22 || size > MaxBytes {
		return nil, fmt.Errorf("archive size must be between 22 bytes and 1 GiB")
	}
	prepared, err := checkDirectory(ctx, source, size)
	if err != nil {
		return nil, fmt.Errorf("validate ZIP directory: %w", err)
	}
	reader, err := zip.NewReader(prepared, prepared.offset+int64(len(prepared.metadata)))
	if err != nil {
		return nil, fmt.Errorf("open X export ZIP: %w", err)
	}
	if len(reader.File) > maxEntries {
		return nil, fmt.Errorf("archive entry count exceeded")
	}
	files := map[string]*zip.File{}
	parts := map[int]*zip.File{}
	var account *zip.File
	for _, file := range reader.File {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("inspect ZIP entries: %w", err)
		}
		name := strings.TrimSuffix(file.Name, "/")
		if name == "" || !utf8.ValidString(name) || strings.ContainsAny(name, "\\:\x00") || strings.HasPrefix(name, "/") || path.Clean(name) != name || name == ".." || strings.HasPrefix(name, "../") {
			return nil, fmt.Errorf("unsafe ZIP entry name")
		}
		if !file.Mode().IsRegular() && !file.Mode().IsDir() {
			return nil, fmt.Errorf("unsupported ZIP entry type")
		}
		if _, exists := files[file.Name]; exists {
			return nil, fmt.Errorf("duplicate ZIP entry")
		}
		files[file.Name] = file
		switch {
		case file.Name == "data/account.js":
			account = file
		case file.Name == "data/tweets.js":
			parts[0] = file
		case strings.HasPrefix(file.Name, "data/tweets-part") && strings.HasSuffix(file.Name, ".js"):
			number := strings.TrimSuffix(strings.TrimPrefix(file.Name, "data/tweets-part"), ".js")
			index, err := strconv.Atoi(number)
			if err != nil || index < 1 || strconv.Itoa(index) != number || index > maxEntries {
				return nil, fmt.Errorf("invalid tweet part name")
			}
			parts[index] = file
		}
	}
	if account == nil || parts[0] == nil {
		return nil, fmt.Errorf("export requires data/account.js and data/tweets.js")
	}
	decoded := int64(0)
	accounts, err := readPart(ctx, account, "account", 0, &decoded)
	if err != nil {
		return nil, fmt.Errorf("decode archive account: %w", err)
	}
	if len(accounts) != 1 {
		return nil, fmt.Errorf("export must contain exactly one account")
	}
	accountRecord, err := parseAccount(accounts[0])
	if err != nil {
		return nil, fmt.Errorf("decode archive account: %w", err)
	}
	id := accountRecord.Id
	if !model.ValidID(id) {
		return nil, fmt.Errorf("archive account has no valid accountId")
	}
	return &Export{AccountID: id, account: accountRecord, files: files, decoded: decoded, records: 1, retained: len(accounts[0]), cache: map[string][]json.RawMessage{"account": accounts}}, nil
}

// checkDirectory validates and captures actual headers before stdlib allocation.
// Ordinary ZIPs require a contiguous directory ending at the selected footer;
// ZIP64, multi-disk layouts and inconsistent directory lengths are unsupported.
func checkDirectory(ctx context.Context, source io.ReaderAt, size int64) (*archiveReader, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("before ZIP directory access: %w", err)
	}
	length := min(size, 22+65535)
	tail := make([]byte, int(length))
	if _, err := source.ReadAt(tail, size-length); err != nil {
		return nil, fmt.Errorf("read ZIP footer: %w", err)
	}
	for index := len(tail) - 22; index >= 0; index-- {
		if !bytes.Equal(tail[index:index+4], []byte{'P', 'K', 5, 6}) {
			continue
		}
		footer := tail[index:]
		if int(binary.LittleEndian.Uint16(footer[20:22]))+22 != len(footer) {
			continue
		}
		count := binary.LittleEndian.Uint16(footer[10:12])
		directorySize := int64(binary.LittleEndian.Uint32(footer[12:16]))
		offset := int64(binary.LittleEndian.Uint32(footer[16:20]))
		end := size - length + int64(index)
		if binary.LittleEndian.Uint16(footer[4:6]) != 0 || binary.LittleEndian.Uint16(footer[6:8]) != 0 || binary.LittleEndian.Uint16(footer[8:10]) != count || count > maxEntries || directorySize > 8<<20 || offset+directorySize != end {
			return nil, fmt.Errorf("unsupported or oversized ZIP directory")
		}
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("before ZIP directory parsing: %w", err)
		}
		metadata := make([]byte, int(directorySize)+22)
		cursor := int64(0)
		for entry := uint16(0); entry < count; entry++ {
			if err := ctx.Err(); err != nil {
				return nil, fmt.Errorf("scan ZIP directory: %w", err)
			}
			if directorySize-cursor < 46 {
				return nil, fmt.Errorf("truncated ZIP directory header")
			}
			header := metadata[cursor : cursor+46]
			if _, err := source.ReadAt(header, offset+cursor); err != nil {
				return nil, fmt.Errorf("read ZIP directory header: %w", err)
			}
			if binary.LittleEndian.Uint32(header) != 0x02014b50 || binary.LittleEndian.Uint16(header[34:36]) != 0 {
				return nil, fmt.Errorf("invalid ZIP directory header")
			}
			recordSize := int64(46) + int64(binary.LittleEndian.Uint16(header[28:30])) + int64(binary.LittleEndian.Uint16(header[30:32])) + int64(binary.LittleEndian.Uint16(header[32:34]))
			if recordSize > directorySize-cursor {
				return nil, fmt.Errorf("ZIP directory record exceeds declared range")
			}
			if err := ctx.Err(); err != nil {
				return nil, fmt.Errorf("read ZIP directory fields: %w", err)
			}
			if recordSize > 46 {
				if _, err := source.ReadAt(metadata[cursor+46:cursor+recordSize], offset+cursor+46); err != nil {
					return nil, fmt.Errorf("read ZIP directory fields: %w", err)
				}
			}
			cursor += recordSize
		}
		if cursor != directorySize {
			return nil, fmt.Errorf("ZIP directory count or size mismatch")
		}
		// The pinned reader also treats size 0xffff as a ZIP64 trigger. A locator
		// in the last record's comment must not redirect it into unchecked file data.
		if directorySize == 0xffff && binary.LittleEndian.Uint32(metadata[directorySize-20:]) == 0x07064b50 {
			return nil, fmt.Errorf("ZIP64 directory locator is unsupported")
		}
		copy(metadata[directorySize:], footer[:22])
		// Drop the unused comment from this view so both parsers see the same footer.
		binary.LittleEndian.PutUint16(metadata[len(metadata)-2:], 0)
		return &archiveReader{ctx: ctx, source: source, offset: offset, metadata: metadata}, nil
	}
	return nil, fmt.Errorf("ZIP footer is missing or unsupported")
}

// readPart decodes a literal YTD array assignment and rejects all trailing code.
func readPart(ctx context.Context, file *zip.File, category string, index int, decoded *int64) ([]json.RawMessage, error) {
	if file.UncompressedSize64 > maxEntryBytes || file.UncompressedSize64 > uint64(maxDecodedBytes-*decoded) {
		return nil, fmt.Errorf("archive decoded-data bound exceeded")
	}
	entry, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open export data file: %w", err)
	}
	data, readErr := io.ReadAll(io.LimitReader(contextReader{ctx: ctx, source: entry}, maxEntryBytes+1))
	closeErr := entry.Close()
	if readErr != nil {
		return nil, fmt.Errorf("read export data file: %w", readErr)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("close export data file: %w", closeErr)
	}
	*decoded += int64(len(data))
	if len(data) > maxEntryBytes || *decoded > maxDecodedBytes {
		return nil, fmt.Errorf("archive decoded-data bound exceeded")
	}
	data = bytes.TrimSpace(data)
	if !utf8.Valid(data) {
		return nil, fmt.Errorf("export data is not valid UTF-8")
	}
	prefix := []byte(fmt.Sprintf("window.YTD.%s.part%d", category, index))
	if !bytes.HasPrefix(data, prefix) {
		return nil, fmt.Errorf("unsupported export data wrapper")
	}
	data = bytes.TrimSpace(data[len(prefix):])
	if len(data) == 0 || data[0] != '=' {
		return nil, fmt.Errorf("invalid export assignment")
	}
	data = bytes.TrimSpace(data[1:])
	if len(data) > 0 && data[len(data)-1] == ';' {
		data = bytes.TrimSpace(data[:len(data)-1])
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("decode export array: %w", err)
	}
	if token != json.Delim('[') {
		return nil, fmt.Errorf("export data must be an array")
	}
	rows := []json.RawMessage{}
	for decoder.More() {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("decode export rows: %w", err)
		}
		if len(rows) >= maxRecords {
			return nil, fmt.Errorf("export row count exceeded")
		}
		var row json.RawMessage
		if err := decoder.Decode(&row); err != nil {
			return nil, fmt.Errorf("decode export row: %w", err)
		}
		if len(row) > maxRecordBytes {
			return nil, fmt.Errorf("export record exceeds byte bound")
		}
		if err := ValidateJSON(row); err != nil {
			return nil, fmt.Errorf("validate export row: %w", err)
		}
		rows = append(rows, row)
	}
	if _, err := decoder.Token(); err != nil {
		return nil, fmt.Errorf("close export array: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("export data contains trailing or executable content")
	}
	return rows, nil
}

// number accepts exact integer JSON values or the export's quoted counters/offsets.
func number(raw json.RawMessage) (int64, error) {
	value := string(raw)
	if len(raw) > 0 && raw[0] == '"' {
		if err := json.Unmarshal(raw, &value); err != nil {
			return 0, fmt.Errorf("decode quoted number: %w", err)
		}
	}
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("decode integer: %w", err)
	}
	if result < 0 {
		return 0, fmt.Errorf("negative counter or offset")
	}
	return result, nil
}
