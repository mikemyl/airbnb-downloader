package airbnb

import (
	"fmt"

	"github.com/go-rod/rod"
)

func (c *Client) GetReviews(page *rod.Page) ([]string, error) {
	searchResults, err := page.Search("div[data-section-id='OVERVIEW_DEFAULT_V2'] a")
	if err != nil {
		return nil, fmt.Errorf("failed to find reviews link: %w", err)
	}
	aLink := searchResults.First
	if err = aLink.Timeout(defaultWaitTime).Click("left", 1); err != nil {
		return nil, fmt.Errorf("failed to click reviews link: %w", err)
	}

	return nil, nil
}
