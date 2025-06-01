package models

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
)

// Document represents the database model for a document
type Document struct {
	BaseModel
	URL            string `gorm:"type:varchar(2048);uniqueIndex"`
	Title          string `gorm:"type:varchar(512)"`
	Content        string `gorm:"type:text"`
	ContentType    string `gorm:"type:varchar(100)"`
	LastCrawled    time.Time
	LastModified   time.Time
	Lang           string `gorm:"type:varchar(10)"`
	MetaDesc       string `gorm:"type:text"`
	ContentLength  int
	ImportanceRank float64
	IndexID        string `gorm:"type:varchar(36);index"`

	DocumentMetadata DocumentMetadata  `gorm:"foreignKey:DocumentID"`
	DocumentLinks    []DocumentLink    `gorm:"foreignKey:SourceID"`
	DocumentKeywords []DocumentKeyword `gorm:"foreignKey:DocumentID"`
}

// BeforeCreate is a GORM hook that generates a UUID if ID is empty
func (d *Document) BeforeCreate() (err error) {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	return
}

// ToDomain converts the database model to a domain entity
func (d *Document) ToDomain() *domain.Document {
	doc := &domain.Document{
		ID:             d.ID,
		URL:            d.URL,
		Title:          d.Title,
		Content:        d.Content,
		ContentType:    domain.ContentType(d.ContentType),
		LastCrawled:    d.LastCrawled,
		LastModified:   d.LastModified,
		Lang:           d.Lang,
		MetaDesc:       d.MetaDesc,
		ContentLength:  d.ContentLength,
		ImportanceRank: d.ImportanceRank,
		IndexID:        d.IndexID,
		StatusCode:     http.StatusOK,
		MetaKeywords:   make([]string, 0),
		Links:          make([]string, 0),
		ParsedContent:  make(map[string]interface{}),
	}

	for _, keyword := range d.DocumentKeywords {
		doc.MetaKeywords = append(doc.MetaKeywords, keyword.Keyword)
	}

	for _, link := range d.DocumentLinks {
		doc.Links = append(doc.Links, link.TargetURL)
	}

	return doc
}
