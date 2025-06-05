package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/mohamedshehata15/intelli-index/internal/adapters/outgoing/storage/models"
	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
	"github.com/mohamedshehata15/intelli-index/internal/core/ports/outgoing"
)

// DocumentRepository implements the outgoing.DocumentRepository interface using GORM
type DocumentRepository struct {
	db *gorm.DB
}

// NewDocumentRepository creates a new document repository
type NewDocumentRepository struct {
	db *gorm.DB
}

// Ensure DocumentRepository implements the outgoing.DocumentRepository interface
var _ outgoing.DocumentRepository = (*DocumentRepository)(nil)

func (d DocumentRepository) Save(ctx context.Context, document *domain.Document) error {
	if document == nil {
		return errors.New("document cannot be nil")
	}

	if err := document.Validate(); err != nil {
		return fmt.Errorf("invalid document: %w", err)
	}

	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		dbDoc := models.Document{}
		dbDoc.FromDomain(document)

		if err := tx.Create(&dbDoc).Error; err != nil {
			if d.isUniqueConstraintViolation(err) {
				return fmt.Errorf("document with URL %s already exists", document.URL)
			}
			return fmt.Errorf("failed to save document: %w", err)
		}

		if document.ID == "" {
			document.ID = dbDoc.ID
		}

		if document.ParsedContent != nil {
			if err := d.saveDocumentMetadata(tx, dbDoc.ID, document.ParsedContent); err != nil {
				return err
			}
		}

		if len(document.MetaKeywords) > 0 {
			if err := d.saveDocumentKeywords(tx, dbDoc.ID, document.MetaKeywords); err != nil {
				return err
			}
		}

		if len(document.Links) > 0 {
			if err := d.saveDocumentLinks(tx, dbDoc.ID, document.Links); err != nil {
				return err
			}
		}

		return nil
	})
}

func (d DocumentRepository) saveDocumentMetadata(tx *gorm.DB, documentID string, parsedContent map[string]interface{}) error {
	metadata := models.DocumentMetadata{
		DocumentID: documentID,
	}

	if author, ok := parsedContent["author"].(string); ok {
		metadata.Author = author
	}
	if publisher, ok := parsedContent["publisher"].(string); ok {
		metadata.Publisher = publisher
	}
	if category, ok := parsedContent["category"].(string); ok {
		metadata.Category = category
	}
	if license, ok := parsedContent["license"].(string); ok {
		metadata.License = license
	}

	if createdDateStr, ok := parsedContent["createdDate"].(string); ok && createdDateStr != "" {
		if parsedTime, err := d.parseDocumentDate(createdDateStr); err == nil {
			metadata.CreatedDate = parsedTime
		}
	}

	if err := tx.Create(&metadata).Error; err != nil {
		return fmt.Errorf("failed to save document metadata: %w", err)
	}
	return nil
}

func (d DocumentRepository) parseDocumentDate(dateStr string) (time.Time, error) {
	if parsedTime, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return parsedTime, nil
	}

	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"02/01/2006",
		"January 2, 2006",
		"2 January 2006",
	}

	for _, format := range formats {
		if parsedTime, err := time.Parse(format, dateStr); err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func (d DocumentRepository) saveDocumentKeywords(tx *gorm.DB, documentID string, keywords []string) error {
	if len(keywords) == 0 {
		return nil
	}

	keywordModels := make([]models.DocumentKeyword, 0, len(keywords))
	for _, keyword := range keywords {
		if keyword != "" {
			keywordModels = append(keywordModels, models.DocumentKeyword{
				DocumentID: documentID,
				Keyword:    keyword,
			})
		}
	}

	if len(keywordModels) > 0 {
		if err := tx.Create(&keywordModels).Error; err != nil {
			return fmt.Errorf("failed to save document keywords: %w", err)
		}
	}

	return nil
}

func (d DocumentRepository) saveDocumentLinks(tx *gorm.DB, documentID string, links []string) error {
	if len(links) == 0 {
		return nil
	}

	linkModels := make([]models.DocumentLink, 0, len(links))
	for _, link := range links {
		if link != "" {
			linkModels = append(linkModels, models.DocumentLink{
				SourceID:  documentID,
				TargetURL: link,
			})
		}
	}

	if len(linkModels) > 0 {
		if err := tx.Create(&linkModels).Error; err != nil {
			return fmt.Errorf("failed to save document links: %w", err)
		}
	}

	return nil
}

func (d DocumentRepository) isUniqueConstraintViolation(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "unique constraint") ||
		strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "UNIQUE constraint failed")
}

func (d DocumentRepository) GetByID(ctx context.Context, id string) (*domain.Document, error) {
	if id == "" {
		return nil, errors.New("document ID cannot be empty")
	}

	var dbDoc models.Document
	result := d.db.WithContext(ctx).
		Preload("DocumentMetadata").
		Preload("DocumentLinks").
		Preload("DocumentKeywords").
		First(&dbDoc, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get document by ID: %w", result.Error)
	}
	return dbDoc.ToDomain(), nil
}

func (d DocumentRepository) GetByURL(ctx context.Context, url string) (*domain.Document, error) {
	//TODO implement me
	panic("implement me")
}

func (d DocumentRepository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (d DocumentRepository) Update(ctx context.Context, document *domain.Document) error {
	//TODO implement me
	panic("implement me")
}

func (d DocumentRepository) List(ctx context.Context, page, pageSize int) ([]*domain.Document, int, error) {
	//TODO implement me
	panic("implement me")
}

func (d DocumentRepository) Search(ctx context.Context, query *domain.SearchQuery) ([]*domain.Document, int, error) {
	//TODO implement me
	panic("implement me")
}

func (d DocumentRepository) CountByIndexID(ctx context.Context, indexID string) (int, error) {
	//TODO implement me
	panic("implement me")
}
