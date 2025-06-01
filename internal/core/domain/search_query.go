package domain

import (
	"errors"
	"strings"
	"time"
)

// SearchQuery represents a search request with filters and pagination
type SearchQuery struct {
	Query               string
	Type                SearchType
	Filters             map[string]interface{}
	Page                int
	PageSize            int
	SortFields          []string
	SortOrder           SortOrder
	IncludeFields       []string
	ExcludeFields       []string
	HighlightFields     []string
	TimeRange           *TimeRange
	Language            string
	FuzzyLevel          int
	FuzzyLevelString    string
	MinimumShouldMatch  string
	ExactTerms          []string
	EntityFilters       map[EntityType][]string
	SearchFields        map[string]float32
	SkipDiversification bool
	UseSearchAfter      bool
	Metadata            map[string]interface{}
}

// SearchType defines the type of search to perform
type SearchType string

const (
	SimpleSearch     SearchType = "simple"
	ExactMatchSearch SearchType = "exact"
	FuzzySearch      SearchType = "fuzzy"
	SemanticSearch   SearchType = "semantic"
)

// SortOrder defines the order of search results
type SortOrder string

const (
	Ascending  SortOrder = "asc"
	Descending SortOrder = "desc"
)

// TimeRange represents a time-based filter
type TimeRange struct {
	Field    string
	From     time.Time
	To       time.Time
	Included bool
}

// NewSearchQuery creates a new search query with default values and validation
func NewSearchQuery(queryText string) (*SearchQuery, error) {
	if strings.TrimSpace(queryText) == "" {
		return nil, errors.New("search query cannot be empty")
	}
	return &SearchQuery{
		Query:              queryText,
		Type:               SimpleSearch,
		Filters:            make(map[string]interface{}),
		EntityFilters:      make(map[EntityType][]string),
		SearchFields:       make(map[string]float32),
		Metadata:           make(map[string]interface{}),
		Page:               1,
		PageSize:           10,
		SortOrder:          Descending,
		FuzzyLevel:         1,
		FuzzyLevelString:   "",
		MinimumShouldMatch: "",
		UseSearchAfter:     false,
	}, nil
}
