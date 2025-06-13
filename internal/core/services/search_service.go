package services

import (
	"context"
	"errors"

	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
	"github.com/mohamedshehata15/intelli-index/internal/core/ports/incoming"
	"github.com/mohamedshehata15/intelli-index/internal/core/ports/outgoing"
)

// searchService implements the incoming.SearchService interface
type searchService struct {
	docRepo   outgoing.DocumentRepository
	indexRepo outgoing.IndexRepository
}

// NewSearchService creates a new search service with the provided dependencies
func NewSearchService(docRepo outgoing.DocumentRepository, indexRepo outgoing.IndexRepository) incoming.SearchService {
	return &searchService{
		docRepo:   docRepo,
		indexRepo: indexRepo,
	}
}

// Ensure searchService implements the incoming.SearchService interface
var _ incoming.SearchService = (*searchService)(nil)

func (s searchService) Search(ctx context.Context, query *domain.SearchQuery) (*domain.SearchResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s searchService) GetDocument(ctx context.Context, id string) (*domain.Document, error) {
	if id == "" {
		return nil, errors.New("document ID cannot be empty")
	}
	return s.docRepo.GetByID(ctx, id)
}

func (s searchService) SuggestQueries(ctx context.Context, partialQuery string, maxSuggestions int) ([]domain.SearchSuggestion, error) {
	//TODO implement me
	panic("implement me")
}

func (s searchService) TrackSuggestionSelection(ctx context.Context, suggestion domain.SearchSuggestion, partialQuery string, position int, timeTakenMs int64) error {
	//TODO implement me
	panic("implement me")
}
