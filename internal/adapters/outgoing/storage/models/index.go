package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Index represents the database model for a search index
type Index struct {
	BaseModel
	Name         string     `gorm:"type:varchar(255);uniqueIndex"`
	Description  string     `gorm:"type:text"`
	SettingsJSON string     `gorm:"type:text;column:settings"`
	MappingsJSON string     `gorm:"type:text;column:mappings"`
	Documents    []Document `gorm:"foreignKey:IndexID"`
}

// BeforeCreate is a GORM hook that generates a UUID if ID is empty
func (i *Index) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == "" {
		i.ID = uuid.NewString()
	}
	return
}
