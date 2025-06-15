package domain

// FileType represents the type of document or file
type FileType string

const (
	FileTypeHTML     FileType = "html"
	FileTypePDF      FileType = "pdf"
	FileTypeText     FileType = "text"
	FileTypeImage    FileType = "image"
	FileTypeVideo    FileType = "video"
	FileTypeAudio    FileType = "audio"
	FileTypeDocument FileType = "document"
	FileTypeArchive  FileType = "archive"
	FileTypeXML      FileType = "xml"
	FileTypeJSON     FileType = "json"
	FileTypeUnknown  FileType = "unknown"
)
