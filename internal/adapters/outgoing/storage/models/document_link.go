package models

import "time"

// DocumentLink represents a link between documents
type DocumentLink struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	SourceID  string `gorm:"type:varchar(36);index"`
	TargetURL string `gorm:"type:varchar(2048);index"`
}
