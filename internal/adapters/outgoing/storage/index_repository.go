package storage

import (
	"github.com/mohamedshehata15/intelli-index/internal/core/ports/outgoing"
	"gorm.io/gorm"
)

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

// Ensure IndexRepository implements the outgoing.IndexRepository interface
var _ outgoing.IndexRepository = (*IndexRepository)(nil)
