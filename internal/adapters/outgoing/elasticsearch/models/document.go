package models

import (
	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
	"time"
)

type Document struct {
	ID                 string                 `json:"id"`
	URL                string                 `json:"url"`
	Title              string                 `json:"title"`
	Content            string                 `json:"content"`
	ContentType        domain.ContentType     `json:"content_type"`
	ContentFingerprint string                 `json:"content_fingerprint"`
	LastCrawled        time.Time              `json:"last_crawled"`
	LastModified       time.Time              `json:"last_modified"`
	Lang               string                 `json:"lang"`
	MetaDesc           string                 `json:"meta_desc"`
	MetaKeywords       []string               `json:"meta_keywords"`
	EnhancedKeywords   []Keyword              `json:"enhanced_keywords"`
	Links              []string               `json:"links"`
	StatusCode         int                    `json:"status_code"`
	ContentLength      int                    `json:"content_length"`
	ImportanceRank     float64                `json:"importance_rank"`
	IndexID            string                 `json:"index_id"`
	IsDuplicate        bool                   `json:"is_duplicate"`
	OriginalDocID      string                 `json:"original_doc_id"`
	VersionCount       int                    `json:"version_count"`
	CurrentVersion     int                    `json:"current_version"`
	ParsedContent      map[string]interface{} `json:"parsed_content"`
	Score              float64                `json:"score"`
}

type Keyword struct {
	Text             string  `json:"text"`
	Score            float64 `json:"score"`
	IsDomainSpecific bool    `json:"is_domain_specific"`
	Category         string  `json:"category"`
	Position         int     `json:"position"`
}
