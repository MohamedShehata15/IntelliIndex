package domain

import (
	"errors"
	"net/url"
	"strings"
)

// URL represents a validated URL as a value object
type URL struct {
	raw          string
	parsed       *url.URL
	normalized   string
	isNormalized bool
}

func NewURL(rawURL string) (*URL, error) {
	if rawURL == "" {
		return nil, errors.New("URL cannot be empty")
	}
	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if parsedURL.Host == "" {
		return nil, errors.New("URL must have a host")
	}
	return &URL{
		raw:    rawURL,
		parsed: parsedURL,
	}, nil
}

func (u *URL) String() string {
	if u.isNormalized {
		return u.normalized
	}
	return u.raw
}

// Normalize returns a normalized version of the URL
func (u *URL) Normalize() string {
	if u.isNormalized {
		return u.normalized
	}

	normalized := *u.parsed
	u.normalizeHost(&normalized)
	u.normalizePath(&normalized)
	u.removeTrackingParams(&normalized)
	normalized.Fragment = ""

	u.normalized = normalized.String()
	u.isNormalized = true

	return u.normalized
}

func (u *URL) normalizeHost(parsedURL *url.URL) {
	parsedURL.Host = strings.ToLower(parsedURL.Host)
	if parsedURL.Port() == "80" && parsedURL.Scheme == "http" {
		parsedURL.Host = parsedURL.Hostname()
	} else if parsedURL.Port() == "443" && parsedURL.Scheme == "https" {
		parsedURL.Host = parsedURL.Hostname()
	}
}

func (u *URL) normalizePath(parsedURL *url.URL) {
	if parsedURL.Path != "/" && strings.HasSuffix(parsedURL.Path, "/") {
		parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")
	}
}

func (u *URL) removeTrackingParams(parsedURL *url.URL) {
	query := parsedURL.Query()
	for _, param := range []string{
		"utm_source", "utm_medium", "utm_campaign", "utm_term", "utm_content", "fbclid", "gclid",
	} {
		query.Del(param)
	}
	parsedURL.RawQuery = query.Encode()
}
