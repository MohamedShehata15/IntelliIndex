package models

import "time"

// DocumentKeyword represents a keyword for a document
type DocumentKeyword struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt  time.Time
	DocumentID string `gorm:"type:varchar(36);index"`
	Keyword    string `gorm:"type:varchar(100);index"`
}
