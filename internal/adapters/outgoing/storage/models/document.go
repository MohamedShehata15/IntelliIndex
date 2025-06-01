package models

import "time"

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
