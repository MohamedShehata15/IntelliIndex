package domain

import "strings"

func init() {
	// Initialize map with all defined content types
	contentTypeMap = map[string]ContentType{
		string(ContentTypeTextPlain):             ContentTypeTextPlain,
		string(ContentTypeTextHTML):              ContentTypeTextHTML,
		string(ContentTypeTextCSS):               ContentTypeTextCSS,
		string(ContentTypeTextJavaScript):        ContentTypeTextJavaScript,
		string(ContentTypeTextXML):               ContentTypeTextXML,
		string(ContentTypeTextCSV):               ContentTypeTextCSV,
		string(ContentTypeTextMarkdown):          ContentTypeTextMarkdown,
		string(ContentTypeApplicationJSON):       ContentTypeApplicationJSON,
		string(ContentTypeApplicationXML):        ContentTypeApplicationXML,
		string(ContentTypeApplicationPDF):        ContentTypeApplicationPDF,
		string(ContentTypeApplicationZip):        ContentTypeApplicationZip,
		string(ContentTypeApplicationForm):       ContentTypeApplicationForm,
		string(ContentTypeApplicationWord):       ContentTypeApplicationWord,
		string(ContentTypeApplicationExcel):      ContentTypeApplicationExcel,
		string(ContentTypeApplicationPowerPoint): ContentTypeApplicationPowerPoint,
		string(ContentTypeImageJPEG):             ContentTypeImageJPEG,
		string(ContentTypeImagePNG):              ContentTypeImagePNG,
		string(ContentTypeImageGIF):              ContentTypeImageGIF,
		string(ContentTypeImageSVG):              ContentTypeImageSVG,
		string(ContentTypeImageWebP):             ContentTypeImageWebP,
		string(ContentTypeAudioMP3):              ContentTypeAudioMP3,
		string(ContentTypeAudioWAV):              ContentTypeAudioWAV,
		string(ContentTypeAudioOGG):              ContentTypeAudioOGG,
		string(ContentTypeVideoMP4):              ContentTypeVideoMP4,
		string(ContentTypeVideoWebM):             ContentTypeVideoWebM,
		string(ContentTypeVideoOGG):              ContentTypeVideoOGG,
		string(ContentTypeOctetStream):           ContentTypeOctetStream,
	}
}

// ContentType represents MIME types content
type ContentType string

const (
	// Text content types
	ContentTypeTextPlain      ContentType = "text/plain"
	ContentTypeTextHTML       ContentType = "text/html"
	ContentTypeTextCSS        ContentType = "text/css"
	ContentTypeTextJavaScript ContentType = "text/javascript"
	ContentTypeTextXML        ContentType = "text/xml"
	ContentTypeTextCSV        ContentType = "text/csv"
	ContentTypeTextMarkdown   ContentType = "text/markdown"

	// Application content types
	ContentTypeApplicationJSON ContentType = "application/json"
	ContentTypeApplicationXML  ContentType = "application/xml"
	ContentTypeApplicationPDF  ContentType = "application/pdf"
	ContentTypeApplicationZip  ContentType = "application/zip"
	ContentTypeApplicationForm ContentType = "application/x-www-form-urlencoded"

	// Document types
	ContentTypeApplicationWord       ContentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	ContentTypeApplicationExcel      ContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	ContentTypeApplicationPowerPoint ContentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"

	// Image content types
	ContentTypeImageJPEG ContentType = "image/jpeg"
	ContentTypeImagePNG  ContentType = "image/png"
	ContentTypeImageGIF  ContentType = "image/gif"
	ContentTypeImageSVG  ContentType = "image/svg+xml"
	ContentTypeImageWebP ContentType = "image/webp"

	// Audio content types
	ContentTypeAudioMP3 ContentType = "audio/mpeg"
	ContentTypeAudioWAV ContentType = "audio/wav"
	ContentTypeAudioOGG ContentType = "audio/ogg"

	// Video content types
	ContentTypeVideoMP4  ContentType = "video/mp4"
	ContentTypeVideoWebM ContentType = "video/webm"
	ContentTypeVideoOGG  ContentType = "video/ogg"

	// Binary type
	ContentTypeOctetStream ContentType = "application/octet-stream"

	// Unknown content type
	ContentTypeUnknown ContentType = ""
)

var contentTypeMap map[string]ContentType

// ParseContentType converts a string MIME type to a ContentType
func ParseContentType(mimeType string) ContentType {
	// Clean up content type by removing parameters
	if idx := strings.Index(mimeType, ";"); idx != -1 {
		mimeType = mimeType[:idx]
	}
	mimeType = strings.TrimSpace(mimeType)

	// Check if it's a known content type
	if contentType, ok := contentTypeMap[mimeType]; ok {
		return contentType
	}

	// Check for general type prefixes
	prefixes := []string{"text/", "image/", "audio/", "video/"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(mimeType, prefix) {
			return ContentType(mimeType)
		}
	}
	return ContentTypeUnknown
}
