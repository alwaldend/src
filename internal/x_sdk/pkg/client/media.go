package client

import "encoding/json"

// AppendMediaUploadRequestMedia holds the JSON encoding of an upload segment.
// JSON requests carry a base64 string encoded as a JSON string literal.
type AppendMediaUploadRequestMedia = json.RawMessage

// ChatMediaUploadAppendRequestMedia holds the JSON encoding of an XChat segment.
type ChatMediaUploadAppendRequestMedia = json.RawMessage

// MediaUploadRequestMedia holds the JSON encoding of a complete media upload.
type MediaUploadRequestMedia = json.RawMessage
