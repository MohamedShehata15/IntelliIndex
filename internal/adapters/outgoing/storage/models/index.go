package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
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

// ToDomain converts the database model to a domain entity
func (i *Index) ToDomain() (*domain.Index, error) {
	index, err := domain.NewIndex(i.Name, i.Description)
	if err != nil {
		return nil, err
	}

	index.ID = i.ID

	if err := i.parseSettings(index); err != nil {
		return nil, err
	}

	if err := i.parseMappings(index); err != nil {
		return nil, err
	}

	return index, nil
}

// parseSettings unmarshal settings JSON into the domain index
func (i *Index) parseSettings(index *domain.Index) error {
	if i.SettingsJSON == "" {
		return nil
	}

	var settings domain.IndexSettings
	if err := json.Unmarshal([]byte(i.SettingsJSON), &settings); err != nil {
		return err
	}
	index.Settings = settings
	return nil
}

// parseMappings unmarshal mapping JSON into the domain index
func (i *Index) parseMappings(index *domain.Index) error {
	if i.MappingsJSON == "" {
		return nil
	}
	var mappings map[string]string
	if err := json.Unmarshal([]byte(i.MappingsJSON), &mappings); err != nil {
		return err
	}
	index.DocumentMapping = mappings
	return nil
}

// FromDomain converts a domain entity to a database model
func (i *Index) FromDomain(index *domain.Index) error {
	i.ID = index.ID
	i.Name = index.Name
	i.Description = index.Description

	settingsJSON, err := json.Marshal(index.Settings)
	if err != nil {
		return err
	}
	i.SettingsJSON = string(settingsJSON)

	mappingsJSON, err := json.Marshal(index.DocumentMapping)
	if err != nil {
		return err
	}
	i.MappingsJSON = string(mappingsJSON)

	return nil
}
