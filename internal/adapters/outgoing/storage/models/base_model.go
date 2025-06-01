package models

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        string `gorm:"type:varchar(36);primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
