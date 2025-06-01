package models

import "time"

// DocumentMetadata represents additional metadata for a document
type DocumentMetadata struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DocumentID string `gorm:"type:varchar(36);uniqueIndex"`
	Author     string `gorm:"type:varchar(255)"`
	Publisher  string `gorm:"type:varchar(255)"`
	Category   string `gorm:"type:varchar(100)"`
	License    string `gorm:"type:varchar(100)"`
}
