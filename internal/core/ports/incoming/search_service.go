package incoming

import (
	"context"
	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
)

// SearchService defines the primary port for the search functionality
type SearchService interface {
	Search(ctx context.Context, query *domain.SearchQuery) (*domain.SearchResult, error)
	GetDocument(ctx context.Context, id string) (*domain.Document, error)
	SuggestQueries(ctx context.Context, partialQuery string, maxSuggestions int) ([]domain.SearchSuggestion, error)
	TrackSuggestionSelection(ctx context.Context, suggestion domain.SearchSuggestion, partialQuery string, position int, timeTakenMs int64) error
}
