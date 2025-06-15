package domain

import "time"

type CrawlResult struct {
	URL               string
	StatusCode        int
	ContentType       string
	Title             string
	Content           string
	Links             []string
	MetaData          map[string]string
	ContentLength     int
	Error             error
	Language          string
	FileType          FileType
	ContentCategory   ContentCategory
	Topics            []string
	Keywords          []Keyword
	ReadingLevel      string
	AuthorInfo        string
	PublishedDate     time.Time
	CleanedContent    string
	WordCount         int
	ImageCount        int
	HasStructuredData bool
	ContentFeatures   map[string]bool
	ClassifierScores  map[string]float64
	IsLowValue        bool
}
