//go:build js && wasm

package browser

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"syscall/js"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/internal/timezone"
	"github.com/spf13/afero"
)

// filesystem lends synchronous Afero operations to one Rust WASI instance.
// Callbacks must use an in-memory filesystem which never waits for JavaScript.
// Sync flushes that filesystem only; durable browser snapshots belong to the app.
type filesystem struct {
	// storage is borrowed from the application and outlives this server instance.
	storage afero.Fs
	// files contains owned open descriptors, including the root preopen at three.
	files map[int]*fileDescriptor
	// nextDescriptor allocates descriptor identities without recycling closed ones.
	nextDescriptor int
	// inodes preserves file identity across opens and renames during this instance.
	inodes map[string]uint64
	// nextInode allocates synthetic inode identities for Afero's portable metadata.
	nextInode uint64
	// location supplies the browser's named local timezone using Go's timezone data.
	location *time.Location
	// endpoint lends the single internal filesystem dispatch callback to JavaScript.
	endpoint js.Value
	// callback remains registered until the Rust backend has been destroyed.
	callback js.Func
	// closeOnce prevents releasing callbacks or descriptors more than once.
	closeOnce sync.Once
	// closeError preserves descriptor cleanup failures for repeated Close calls.
	closeError error
}

// fileDescriptor retains access rights and metadata missing from afero.File.
type fileDescriptor struct {
	// file is closed by fd_close or the filesystem owner after guest termination.
	file afero.File
	// name follows renames so relative directory operations resolve consistently.
	name string
	// inode remains stable even after an open file is renamed or unlinked.
	inode uint64
	// directory distinguishes directories from ordinary byte-addressable files.
	directory bool
	// readable records the read capability requested by the guest.
	readable bool
	// writable records the write capability requested by the guest.
	writable bool
	// flags retains WASI append, sync and nonblocking descriptor flags.
	flags int
	// entries holds one directory enumeration until the guest restarts at cookie zero.
	entries []any
}

// newFilesystem prepares a root preopen and writable temporary storage without
// replacing the caller's filesystem. Its descriptors and callback are owned.
func newFilesystem(storage afero.Fs) (*filesystem, error) {
	if storage == nil {
		return nil, fmt.Errorf("prepare browser Anki filesystem: no filesystem supplied")
	}
	zone := js.Global().Get("Intl").Get("DateTimeFormat").New().Call("resolvedOptions").Get("timeZone").String()
	location, err := time.LoadLocation(zone)
	if err != nil {
		return nil, fmt.Errorf("load browser timezone %q: %w", zone, err)
	}
	if err := storage.MkdirAll("/tmp", 0o700); err != nil {
		return nil, fmt.Errorf("prepare browser Anki temporary directory: %w", err)
	}
	root, err := storage.Open("/")
	if err != nil {
		return nil, fmt.Errorf("open browser Anki filesystem root: %w", err)
	}
	fs := &filesystem{
		storage: storage, files: make(map[int]*fileDescriptor), nextDescriptor: 4,
		inodes: map[string]uint64{"/": 1}, nextInode: 2, location: location,
	}
	fs.files[3] = &fileDescriptor{file: root, name: "/", inode: 1, directory: true, readable: true}
	fs.callback = js.FuncOf(fs.invoke)
	fs.endpoint = js.Global().Get("Object").New()
	fs.endpoint.Set("call", fs.callback)
	return fs, nil
}

// Close releases every remaining descriptor and the callback after guest execution.
func (fs *filesystem) Close() error {
	fs.closeOnce.Do(func() {
		var failures []error
		for number, descriptor := range fs.files {
			if err := descriptor.file.Close(); err != nil {
				failures = append(failures, fmt.Errorf("close browser descriptor %d (%q): %w", number, descriptor.name, err))
			}
			delete(fs.files, number)
		}
		fs.closeError = errors.Join(failures...)
		fs.callback.Release()
	})
	return fs.closeError
}

// invoke translates filesystem failures into WASI errno without throwing across
// Wasm runtimes. wasi.js owns ABI validation and supplies the operation arguments.
func (fs *filesystem) invoke(_ js.Value, arguments []js.Value) any {
	if len(arguments) == 0 || arguments[0].Type() != js.TypeString {
		return map[string]any{"errno": 28}
	}
	value, err := fs.operate(arguments[0].String(), arguments[1:])
	if err != nil {
		return map[string]any{"errno": filesystemErrno(err)}
	}
	return map[string]any{"errno": 0, "value": value}
}

// operate implements the private operations used by the WASI memory codec.
func (fs *filesystem) operate(operation string, arguments []js.Value) (any, error) {
	switch operation {
	case "open":
		return fs.open(arguments[0].Int(), arguments[1].String(), arguments[2].Int(), arguments[3].Bool(), arguments[4].Bool(), arguments[5].Int())
	case "pathStat", "mkdir", "removeDirectory", "unlink", "pathTimes":
		name, err := fs.resolve(arguments[0].Int(), arguments[1].String())
		if err != nil {
			return nil, fmt.Errorf("resolve browser filesystem path: %w", err)
		}
		switch operation {
		case "pathStat":
			info, err := fs.storage.Stat(name)
			if err != nil {
				return nil, fmt.Errorf("stat browser file %q: %w", name, err)
			}
			return fs.stat(info, fs.inode(name)), nil
		case "mkdir":
			if err := fs.parent(name); err != nil {
				return nil, fmt.Errorf("create browser directory %q: %w", name, err)
			}
			if err := fs.storage.Mkdir(name, 0o700); err != nil {
				return nil, fmt.Errorf("create browser directory %q: %w", name, err)
			}
		case "removeDirectory", "unlink":
			if err := fs.remove(name, operation == "removeDirectory"); err != nil {
				return nil, fmt.Errorf("remove browser path %q: %w", name, err)
			}
		case "pathTimes":
			if err := fs.setTimes(name, arguments[2].String(), arguments[3].String(), arguments[4].Int()); err != nil {
				return nil, fmt.Errorf("set browser file times %q: %w", name, err)
			}
		}
		return nil, nil
	case "rename":
		oldName, oldError := fs.resolve(arguments[0].Int(), arguments[1].String())
		newName, newError := fs.resolve(arguments[2].Int(), arguments[3].String())
		if err := errors.Join(oldError, newError); err != nil {
			return nil, fmt.Errorf("resolve browser rename paths: %w", err)
		}
		if err := fs.rename(oldName, newName); err != nil {
			return nil, fmt.Errorf("rename browser path %q to %q: %w", oldName, newName, err)
		}
		return nil, nil
	case "timezoneUTC", "timezoneLocal":
		seconds, err := strconv.ParseInt(arguments[0].String(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse browser timezone timestamp: %w", err)
		}
		if operation == "timezoneUTC" {
			return int(timezone.OffsetForUTC(fs.location, seconds)), nil
		}
		return strconv.FormatUint(timezone.OffsetsForLocal(fs.location, seconds), 10), nil
	}
	if len(arguments) == 0 {
		return nil, syscall.EINVAL
	}
	number := arguments[0].Int()
	descriptor, found := fs.files[number]
	if !found {
		return nil, syscall.EBADF
	}
	value, err := fs.descriptorOperation(operation, number, descriptor, arguments[1:])
	if err != nil {
		return nil, fmt.Errorf("browser descriptor %d (%q) %s: %w", number, descriptor.name, operation, err)
	}
	return value, nil
}

// descriptorOperation applies one validated operation to an existing descriptor.
func (fs *filesystem) descriptorOperation(operation string, number int, descriptor *fileDescriptor, arguments []js.Value) (any, error) {
	switch operation {
	case "close":
		delete(fs.files, number)
		if err := descriptor.file.Close(); err != nil {
			return nil, fmt.Errorf("close file: %w", err)
		}
	case "fdstat":
		return map[string]any{"type": fileType(descriptor.directory), "flags": descriptor.flags, "readable": descriptor.readable, "writable": descriptor.writable}, nil
	case "stat":
		info, err := descriptor.file.Stat()
		if err != nil {
			return nil, fmt.Errorf("stat open file: %w", err)
		}
		return fs.stat(info, descriptor.inode), nil
	case "flags":
		descriptor.flags = arguments[0].Int()
	case "sync":
		if err := descriptor.file.Sync(); err != nil {
			return nil, fmt.Errorf("sync working filesystem: %w", err)
		}
	case "seek":
		if descriptor.directory {
			return nil, syscall.EISDIR
		}
		offset, err := strconv.ParseInt(arguments[0].String(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse seek offset: %w", errors.Join(syscall.EINVAL, err))
		}
		whence := arguments[1].Int()
		position, err := descriptor.file.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("read file position: %w", err)
		}
		if whence == io.SeekEnd {
			info, err := descriptor.file.Stat()
			if err != nil {
				return nil, fmt.Errorf("stat seek target: %w", err)
			}
			position = info.Size()
		} else if whence == io.SeekStart {
			position = 0
		} else if whence != io.SeekCurrent {
			return nil, syscall.EINVAL
		}
		if offset < -position || (offset > 0 && position+offset < position) {
			return nil, syscall.EINVAL
		}
		position, err = descriptor.file.Seek(position+offset, io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("seek file: %w", err)
		}
		return strconv.FormatInt(position, 10), nil
	case "read":
		if !descriptor.readable {
			return nil, syscall.EBADF
		}
		if descriptor.directory {
			return nil, syscall.EISDIR
		}
		data := make([]byte, arguments[0].Int())
		count, err := descriptor.file.Read(data)
		// MemMapFs reports ErrUnexpectedEOF after seeking past the current end.
		// WASI regular-file reads still return the available bytes, including zero.
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, fmt.Errorf("read file: %w", err)
		}
		value := js.Global().Get("Uint8Array").New(count)
		js.CopyBytesToJS(value, data[:count])
		return value, nil
	case "write":
		if !descriptor.writable {
			return nil, syscall.EBADF
		}
		if descriptor.directory {
			return nil, syscall.EISDIR
		}
		if descriptor.flags&1 != 0 {
			if _, err := descriptor.file.Seek(0, io.SeekEnd); err != nil {
				return nil, fmt.Errorf("seek append position: %w", err)
			}
		}
		data := make([]byte, arguments[0].Length())
		js.CopyBytesToGo(data, arguments[0])
		count, err := descriptor.file.Write(data)
		if err != nil {
			return nil, fmt.Errorf("write file: %w", err)
		}
		if descriptor.flags&(2|8|16) != 0 {
			if err := descriptor.file.Sync(); err != nil {
				return nil, fmt.Errorf("sync file after write: %w", err)
			}
		}
		return count, nil
	case "truncate":
		if !descriptor.writable {
			return nil, syscall.EBADF
		}
		if descriptor.directory {
			return nil, syscall.EISDIR
		}
		size, err := strconv.ParseInt(arguments[0].String(), 10, 64)
		if err != nil || size < 0 {
			return nil, fmt.Errorf("parse file size: %w", errors.Join(syscall.EINVAL, err))
		}
		if err := descriptor.file.Truncate(size); err != nil {
			return nil, fmt.Errorf("truncate file: %w", err)
		}
	case "times":
		if err := fs.setTimes(descriptor.name, arguments[0].String(), arguments[1].String(), arguments[2].Int()); err != nil {
			return nil, fmt.Errorf("set open file times: %w", err)
		}
	case "readdir":
		if !descriptor.directory {
			return nil, syscall.ENOTDIR
		}
		if descriptor.entries == nil || arguments[0].Bool() {
			entries, err := afero.ReadDir(fs.storage, descriptor.name)
			if err != nil {
				return nil, fmt.Errorf("read directory: %w", err)
			}
			descriptor.entries = make([]any, 0, len(entries))
			for _, entry := range entries {
				descriptor.entries = append(descriptor.entries, map[string]any{"name": entry.Name(), "type": fileType(entry.IsDir()), "inode": strconv.FormatUint(fs.inode(path.Join(descriptor.name, entry.Name())), 10)})
			}
		}
		return descriptor.entries, nil
	default:
		return nil, syscall.ENOSYS
	}
	return nil, nil
}

// resolve accepts relative WASI paths beneath a currently open directory.
func (fs *filesystem) resolve(number int, name string) (string, error) {
	descriptor, found := fs.files[number]
	if !found {
		return "", syscall.EBADF
	}
	if !descriptor.directory {
		return "", syscall.ENOTDIR
	}
	if name == "" {
		return "", syscall.ENOENT
	}
	if strings.ContainsRune(name, 0) || path.IsAbs(name) {
		return "", syscall.EINVAL
	}
	return path.Join(descriptor.name, name), nil
}

// parent enforces existing directory parents even for permissive memory backends.
func (fs *filesystem) parent(name string) error {
	info, err := fs.storage.Stat(path.Dir(name))
	if err != nil {
		return fmt.Errorf("stat parent of %q: %w", name, err)
	}
	if !info.IsDir() {
		return syscall.ENOTDIR
	}
	return nil
}

// open translates WASI creation flags and read/write rights to Afero.
func (fs *filesystem) open(directory int, relative string, flags int, readable, writable bool, descriptorFlags int) (any, error) {
	name, err := fs.resolve(directory, relative)
	if err != nil {
		return nil, fmt.Errorf("resolve open path: %w", err)
	}
	if err := fs.parent(name); err != nil {
		return nil, fmt.Errorf("open browser file %q: %w", name, err)
	}
	info, statError := fs.storage.Stat(name)
	if statError != nil && !errors.Is(statError, os.ErrNotExist) {
		return nil, fmt.Errorf("inspect browser open path %q: %w", name, statError)
	}
	if statError == nil && flags&1 != 0 && flags&4 != 0 {
		return nil, syscall.EEXIST
	}
	if flags&2 != 0 && (statError != nil || !info.IsDir()) {
		if statError != nil {
			return nil, fmt.Errorf("open browser directory %q: %w", name, statError)
		}
		return nil, syscall.ENOTDIR
	}
	if statError == nil && info.IsDir() && (writable || flags&(1|8) != 0) {
		return nil, syscall.EISDIR
	}
	mode := os.O_RDONLY
	if writable && readable {
		mode = os.O_RDWR
	} else if writable {
		mode = os.O_WRONLY
	}
	if flags&1 != 0 {
		mode |= os.O_CREATE
		if flags&4 != 0 {
			mode |= os.O_EXCL
		}
	}
	if flags&8 != 0 {
		if !writable {
			return nil, syscall.EACCES
		}
		mode |= os.O_TRUNC
	}
	file, err := fs.storage.OpenFile(name, mode, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open browser file %q: %w", name, err)
	}
	number := fs.nextDescriptor
	fs.nextDescriptor++
	fs.files[number] = &fileDescriptor{file: file, name: name, inode: fs.inode(name), directory: statError == nil && info.IsDir(), readable: readable, writable: writable, flags: descriptorFlags}
	return number, nil
}

// remove enforces POSIX file/directory distinctions and empty-directory removal.
func (fs *filesystem) remove(name string, directory bool) error {
	info, err := fs.storage.Stat(name)
	if err != nil {
		return fmt.Errorf("inspect removal target: %w", err)
	}
	if directory != info.IsDir() {
		if directory {
			return syscall.ENOTDIR
		}
		return syscall.EISDIR
	}
	if directory {
		entries, err := afero.ReadDir(fs.storage, name)
		if err != nil {
			return fmt.Errorf("inspect directory contents: %w", err)
		}
		if len(entries) != 0 {
			return syscall.ENOTEMPTY
		}
		if name == "/" {
			return syscall.EPERM
		}
	}
	if err := fs.storage.Remove(name); err != nil {
		return fmt.Errorf("remove path: %w", err)
	}
	delete(fs.inodes, name)
	return nil
}

// rename preserves replacement semantics and identities in the shared namespace.
func (fs *filesystem) rename(oldName, newName string) error {
	info, err := fs.storage.Stat(oldName)
	if err != nil {
		return fmt.Errorf("inspect rename source: %w", err)
	}
	if oldName == newName {
		return nil
	}
	if oldName == "/" || newName == "/" || (info.IsDir() && strings.HasPrefix(newName, oldName+"/")) {
		return syscall.EINVAL
	}
	if err := fs.parent(newName); err != nil {
		return fmt.Errorf("inspect rename destination parent: %w", err)
	}
	target, err := fs.storage.Stat(newName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect rename destination: %w", err)
	}
	if err == nil {
		if target.IsDir() != info.IsDir() {
			if target.IsDir() {
				return syscall.EISDIR
			}
			return syscall.ENOTDIR
		}
		if target.IsDir() {
			entries, err := afero.ReadDir(fs.storage, newName)
			if err != nil {
				return fmt.Errorf("inspect rename destination directory: %w", err)
			}
			if len(entries) != 0 {
				return syscall.ENOTEMPTY
			}
		}
	}
	if err := fs.storage.Rename(oldName, newName); err != nil {
		return fmt.Errorf("replace rename destination: %w", err)
	}
	delete(fs.inodes, newName)
	for name, inode := range fs.inodes {
		if name == oldName || strings.HasPrefix(name, oldName+"/") {
			fs.inodes[newName+strings.TrimPrefix(name, oldName)] = inode
			delete(fs.inodes, name)
		}
	}
	for _, descriptor := range fs.files {
		if descriptor.name == oldName || strings.HasPrefix(descriptor.name, oldName+"/") {
			descriptor.name = newName + strings.TrimPrefix(descriptor.name, oldName)
		}
	}
	return nil
}

// setTimes applies explicit/now WASI timestamps; Afero exposes modification time.
func (fs *filesystem) setTimes(name, access, modification string, flags int) error {
	info, err := fs.storage.Stat(name)
	if err != nil {
		return fmt.Errorf("inspect timestamp target: %w", err)
	}
	atime, mtime := info.ModTime(), info.ModTime()
	for _, target := range []struct {
		// text is a decimal nanosecond timestamp encoded without JavaScript rounding.
		text string
		// explicit and now identify the corresponding WASI timestamp flags.
		explicit, now int
		// value receives the selected timestamp.
		value *time.Time
	}{{access, 1, 2, &atime}, {modification, 4, 8, &mtime}} {
		if flags&target.explicit != 0 {
			nanoseconds, err := strconv.ParseUint(target.text, 10, 64)
			if err != nil {
				return fmt.Errorf("parse file timestamp: %w", errors.Join(syscall.EINVAL, err))
			}
			*target.value = time.Unix(int64(nanoseconds/1_000_000_000), int64(nanoseconds%1_000_000_000))
		} else if flags&target.now != 0 {
			*target.value = time.Now()
		}
	}
	if err := fs.storage.Chtimes(name, atime, mtime); err != nil {
		return fmt.Errorf("change file timestamps: %w", err)
	}
	return nil
}

// inode allocates a stable portable file identity on its first observation.
func (fs *filesystem) inode(name string) uint64 {
	if number, found := fs.inodes[name]; found {
		return number
	}
	number := fs.nextInode
	fs.nextInode++
	fs.inodes[name] = number
	return number
}

// stat encodes wide metadata as decimal strings for exact JavaScript conversion.
func (fs *filesystem) stat(info os.FileInfo, inode uint64) any {
	return map[string]any{"type": fileType(info.IsDir()), "inode": strconv.FormatUint(inode, 10), "size": strconv.FormatInt(info.Size(), 10), "time": strconv.FormatInt(info.ModTime().UnixNano(), 10)}
}

// fileType maps the supported Afero file kinds to WASI Preview 1 constants.
func fileType(directory bool) int {
	if directory {
		return 3
	}
	return 4
}

// filesystemErrno translates portable Go errors to the WASI Preview 1 ABI.
func filesystemErrno(err error) int {
	for _, candidate := range []struct {
		// cause identifies the underlying portable filesystem failure.
		cause error
		// number is its WASI errno representation.
		number int
	}{
		{syscall.EBADF, 8},
		{os.ErrClosed, 8},
		{os.ErrExist, 20},
		{syscall.EFBIG, 22},
		{syscall.EINVAL, 28},
		{os.ErrInvalid, 28},
		{syscall.EISDIR, 31},
		{syscall.ELOOP, 32},
		{syscall.EMFILE, 33},
		{syscall.ENAMETOOLONG, 37},
		{os.ErrNotExist, 44},
		{syscall.ENOMEM, 48},
		{syscall.ENOSPC, 51},
		{syscall.ENOSYS, 52},
		{syscall.ENOTDIR, 54},
		{syscall.ENOTEMPTY, 55},
		{syscall.ENOTSUP, 58},
		{syscall.EOVERFLOW, 61},
		{syscall.EPERM, 63},
		{os.ErrPermission, 2},
	} {
		if errors.Is(err, candidate.cause) {
			return candidate.number
		}
	}
	return 29
}
