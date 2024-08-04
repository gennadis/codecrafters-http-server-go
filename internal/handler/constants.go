package handler

// HTTP Status Codes
const (
	StatusOK                  = "200 OK"
	StatusCreated             = "201 Created"
	StatusNotFound            = "404 Not Found"
	StatusInternalServerError = "500 Internal Server Error"
	StatusMethodNotAllowed    = "405 Method Not Allowed"
)

// Headers
const (
	HeaderContentType   = "Content-Type"
	HeaderContentLength = "Content-Length"
	HeaderUserAgent     = "User-Agent"
)

// Content Types
const (
	ContentTypeTextPlain              = "text/plain"
	ContentTypeApplicationOctetStream = "application/octet-stream"
)

// Methods
const (
	MethodGet  = "GET"
	MethodPost = "POST"
)
