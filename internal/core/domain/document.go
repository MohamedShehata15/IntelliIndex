package domain

import (
	"errors"
	"net/http"
	"time"
)

// Document represents a searchable document
type Document struct {
	ID                 string
	URL                string
	Title              string
	Content            string
	ContentType        ContentType
	ContentFingerprint string
	LastCrawled        time.Time
	LastModified       time.Time
	Lang               string
	MetaDesc           string
	MetaKeywords       []string
	EnhancedKeywords   []Keyword
	Links              []string
	StatusCode         int
	ContentLength      int
	ImportanceRank     float64
	IndexID            string
	IsDuplicate        bool
	OriginalDocID      string
	VersionCount       int
	CurrentVersion     int
	ParsedContent      map[string]interface{}
	Score              float64
}

// Keyword represents a document keyword with relevance information
type Keyword struct {
	Text             string
	Score            float64
	IsDomainSpecific bool
	Category         string
	Position         int
}

// NewDocument creates a new Document
func NewDocument(url, title, content, contentType string) (*Document, error) {
	if url == "" {
		return nil, errors.New("URL cannot be empty")
	}
	parsedURL, err := NewURL(url)
	if err != nil {
		return nil, errors.New("invalid URL: " + err.Error())
	}
	normalizedURL := parsedURL.Normalize()
	return &Document{
		URL:            normalizedURL,
		Title:          title,
		Content:        content,
		ContentType:    ParseContentType(contentType),
		LastCrawled:    time.Now(),
		StatusCode:     http.StatusOK,
		MetaKeywords:   make([]string, 0),
		Links:          make([]string, 0),
		VersionCount:   1,
		CurrentVersion: 1,
		IsDuplicate:    false,
		ParsedContent:  make(map[string]interface{}),
	}, nil
}

// UpdateContent updates the document's content & version information
func (d *Document) UpdateContent(content string, contentFingerprint string) {
	if d.Content == content {
		return
	}
	d.Content = content
	d.ContentFingerprint = contentFingerprint
	d.LastModified = time.Now()
	d.VersionCount++
	d.CurrentVersion = d.VersionCount
}

// MarkAsDuplicate marks this document as a duplicate of another document
func (d *Document) MarkAsDuplicate(originalDocID string) {
	d.IsDuplicate = true
	d.OriginalDocID = originalDocID
}

// SetContentFingerprint sets the content fingerprint for this document
func (d *Document) SetContentFingerprint(fingerprint string) {
	d.ContentFingerprint = fingerprint
}

// Validate ensures the document is valid
func (d *Document) Validate() error {
	if d.URL == "" {
		return errors.New("document URL cannot be empty")
	}
	return nil
}
