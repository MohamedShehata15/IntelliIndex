package domain

// FilterRule defines a filtering rule for content
type FilterRule struct {
	ID          string
	Name        string
	Description string
	Type        FilterRuleType
}

// FilterRuleType indicates the type of filter rule
type FilterRuleType string

const (
	// RegexRule applies a regular expression pattern
	RegexRule FilterRuleType = "regex"
	// ContainsRule checks if content contains a string
	ContainsRule FilterRuleType = "contains"
	// ExactMatchRule requires an exact match
	ExactMatchRule FilterRuleType = "exact"
	// DomainRule applies to specific domains
	DomainRule FilterRuleType = "domain"
	// PathRule applies to URL paths
	PathRule FilterRuleType = "path"
	// ContentSizeRule filters based on content size
	ContentSizeRule FilterRuleType = "size"
	// ContentTypeRule filters based on content type
	ContentTypeRule FilterRuleType = "content_type"
)
