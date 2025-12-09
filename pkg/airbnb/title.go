package airbnb

import (
	"fmt"

	"github.com/go-rod/rod"
)

// GetTitle extracts the title from an Airbnb listing page.
func (c *Client) GetTitle(page *rod.Page) (string, error) {
	searchResults, err := page.Search("div[data-section-id='TITLE_DEFAULT'] h1")
	if err != nil {
		return "", fmt.Errorf("failed to find title: %w", err)
	}
	text, err := searchResults.First.Text()
	if err != nil {
		return "", fmt.Errorf("failed to get title text: %w", err)
	}
	return text, nil
}
