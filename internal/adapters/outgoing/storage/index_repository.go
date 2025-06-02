package storage

import "gorm.io/gorm"

// IndexRepository implements the outgoing.IndexRepository interface using GORM
type IndexRepository struct {
	db *gorm.DB
}

// NewIndexRepository creates a new gorm index repository
func NewIndexRepository(db *gorm.DB) *IndexRepository {
	return &IndexRepository{
		db: db,
	}
}
