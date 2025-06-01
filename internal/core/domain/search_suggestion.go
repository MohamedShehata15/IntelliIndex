package domain

import "time"

type SearchSuggestion struct {
	Text          string
	Source        SuggestionSource
	Score         float64
	LastUsedTime  time.Time
	UseCount      int
	CorrectedFrom string
}

// SuggestionSource represents a source type for search suggestions
type SuggestionSource string

const (
	SuggestionSourceHistory      SuggestionSource = "history"
	SuggestionSourcePopular      SuggestionSource = "popular"
	SuggestionSourceAutocomplete SuggestionSource = "autocomplete"
	SuggestionSourceRelated      SuggestionSource = "related"
)
