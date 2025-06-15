package domain

// ContentCategory represents a document content category
type ContentCategory string

const (
	CategoryTechnology    ContentCategory = "technology"
	CategoryBusiness      ContentCategory = "business"
	CategoryHealth        ContentCategory = "health"
	CategoryEducation     ContentCategory = "education"
	CategoryEntertainment ContentCategory = "entertainment"
	CategorySports        ContentCategory = "sports"
	CategoryScience       ContentCategory = "science"
	CategoryNews          ContentCategory = "news"
	CategoryTravel        ContentCategory = "travel"
	CategoryFinance       ContentCategory = "finance"
	CategoryOther         ContentCategory = "other"
)
