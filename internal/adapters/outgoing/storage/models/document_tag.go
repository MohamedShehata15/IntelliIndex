package models

import "time"

// DocumentTag represents a tag for a document
type DocumentTag struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt  time.Time
	DocumentID string `gorm:"type:varchar(36);index"`
	Tag        string `gorm:"type:varchar(100);index"`
}
