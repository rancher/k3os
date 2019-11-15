package upgrade

import (
	"fmt"
	"net/url"
)

// NewChannel creates a new channel for the provided URL
func NewChannel(uri string) (Channel, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if url.Scheme != "github-releases" {
		return nil, fmt.Errorf("scheme not supported: %s", url)
	}
	repo := url.Opaque
	if repo == "" {
		repo = url.Host + url.Path
	}
	if repo == "" {
		return nil, fmt.Errorf("url not supported: %s", url)
	}
	return &github{
		repo: repo,
	}, nil
}
