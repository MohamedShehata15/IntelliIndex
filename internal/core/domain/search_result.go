package domain

// SearchResult represents the result of a search query
type SearchResult struct {
	TotalHits    int
	Documents    []*Document
	Page         int
	PageSize     int
	TotalPages   int
	Took         int64
	Suggestions  []string
	Highlighting map[string]map[string][]string
	QueryID      string
}
