package storage

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm/clause"
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
	if url == "" {
		return nil, errors.New("document URL cannot be empty")
	}

	var dbDoc models.Document
	result := d.db.WithContext(ctx).
		Preload("DocumentMetadata").
		Preload("DocumentLinks").
		Preload("DocumentKeywords").
		First(&dbDoc, "url = ?", url)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get document by URL: ", result.Error)
	}
	return dbDoc.ToDomain(), nil
}

func (d DocumentRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("document ID cannot be empty")
	}

	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exists bool
		if err := tx.Model(&models.Document{}).Select("1").Where("id = ?", id).First(&exists).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("failed to check document existence: %w", err)
		}

		if err := tx.Where("document_id = ?", id).Delete(&models.DocumentKeyword{}).Error; err != nil {
			return fmt.Errorf("failed to delete document keywords: %w", err)
		}

		if err := tx.Where("source_id = ?", id).Delete(&models.DocumentLink{}).Error; err != nil {
			return fmt.Errorf("failed to delete document links: %w", err)
		}

		if err := tx.Where("document_id = ?", id).Delete(&models.DocumentMetadata{}).Error; err != nil {
			return fmt.Errorf("failed to delete document metadata: %w", err)
		}

		if err := tx.Delete(&models.Document{}, "id = ?", id).Error; err != nil {
			return fmt.Errorf("failed to delete document: %w", err)
		}

		return nil
	})
}

func (d DocumentRepository) Update(ctx context.Context, document *domain.Document) error {
	if document == nil {
		return errors.New("document cannot be nil")
	}
	if document.ID == "" {
		return errors.New("document ID cannot be empty")
	}
	if err := document.Validate(); err != nil {
		return fmt.Errorf("invalid document: %w", err)
	}

	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingDoc models.Document
		if err := tx.First(&existingDoc, "id = ?", document.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("document with ID %s does not exist", document.ID)
			}
			return fmt.Errorf("failed to check if document exists: %w", err)
		}

		dbDoc := models.Document{}
		dbDoc.FromDomain(document)
		dbDoc.ID = document.ID

		if err := tx.Model(&dbDoc).Updates(map[string]interface{}{
			"url":             dbDoc.URL,
			"title":           dbDoc.Title,
			"content":         dbDoc.Content,
			"content_type":    dbDoc.ContentType,
			"last_crawled":    dbDoc.LastCrawled,
			"last_modified":   dbDoc.LastModified,
			"lang":            dbDoc.Lang,
			"meta_desc":       dbDoc.MetaDesc,
			"content_length":  dbDoc.ContentLength,
			"importance_rank": dbDoc.ImportanceRank,
			"index_id":        dbDoc.IndexID,
		}).Error; err != nil {
			if d.isUniqueConstraintViolation(err) {
				return fmt.Errorf("document with URL %s already exists", document.URL)
			}
			return fmt.Errorf("failed to update document: %w", err)
		}

		if document.ParsedContent != nil {
			if err := d.updateDocumentMetadata(tx, document.ID, document.ParsedContent); err != nil {
				return err
			}
		}

		if err := d.updateDocumentKeywords(tx, document.ID, document.MetaKeywords); err != nil {
			return err
		}

		if err := d.updateDocumentLinks(tx, document.ID, document.Links); err != nil {
			return err
		}

		return nil
	})
}

func (d DocumentRepository) updateDocumentMetadata(tx *gorm.DB, documentID string, parsedContent map[string]interface{}) error {
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

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "document_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"author", "publisher", "category", "license", "created_date"}),
	}).Create(&metadata).Error; err != nil {
		return fmt.Errorf("failed to update document metadata: %w", err)
	}

	return nil
}

func (d DocumentRepository) updateDocumentKeywords(tx *gorm.DB, documentID string, keywords []string) error {
	// Delete existing keywords
	if err := tx.Where("document_id = ?", documentID).Delete(&models.DocumentKeyword{}).Error; err != nil {
		return fmt.Errorf("failed to delete old document keywords: %w", err)
	}

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

func (d DocumentRepository) updateDocumentLinks(tx *gorm.DB, documentID string, links []string) error {
	if err := tx.Where("source_id = ?", documentID).Delete(&models.DocumentLink{}).Error; err != nil {
		return fmt.Errorf("failed to delete old document links: %w", err)
	}

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

func (d DocumentRepository) List(ctx context.Context, page, pageSize int) ([]*domain.Document, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var count int64
	if err := d.db.WithContext(ctx).Model(&models.Document{}).Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	var dbDocs []models.Document
	result := d.db.WithContext(ctx).
		Preload("DocumentMetadata").
		Preload("DocumentLinks").
		Preload("DocumentKeywords").
		Offset(offset).
		Limit(pageSize).
		Order("last_crawled desc").
		Find(&dbDocs)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to list documents: %w", result.Error)
	}
	documents := make([]*domain.Document, len(dbDocs))
	for i, dbDoc := range dbDocs {
		documents[i] = dbDoc.ToDomain()
	}
	return documents, int(count), nil
}

func (d DocumentRepository) Search(ctx context.Context, query *domain.SearchQuery) ([]*domain.Document, int, error) {
	//TODO implement me
	panic("implement me")
}

func (d DocumentRepository) CountByIndexID(ctx context.Context, indexID string) (int, error) {
	//TODO implement me
	panic("implement me")
}
