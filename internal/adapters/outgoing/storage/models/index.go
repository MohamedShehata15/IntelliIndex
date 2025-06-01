package models

// Index represents the database model for a search index
type Index struct {
	BaseModel
	Name         string     `gorm:"type:varchar(255);uniqueIndex"`
	Description  string     `gorm:"type:text"`
	SettingsJSON string     `gorm:"type:text;column:settings"`
	MappingsJSON string     `gorm:"type:text;column:mappings"`
	Documents    []Document `gorm:"foreignKey:IndexID"`
}
