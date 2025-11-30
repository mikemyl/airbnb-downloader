package airbnb

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// GetListing extracts all information from an Airbnb listing page.
func (c *Client) GetListing(listingURL string) (*Listing, error) {
	parsedURL, err := url.Parse(listingURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse listing url: %w", err)
	}

	target := proto.TargetCreateTarget{URL: listingURL}
	page, err := c.browser.Page(target)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}
	defer func(page *rod.Page) {
		_ = page.Close()
	}(page)

	if err = page.WaitLoad(); err != nil {
		return nil, err
	}

	if !c.hasGonePastTheTheTranslationDialog {
		_ = closeTranslationOnDialog(page)
		c.hasGonePastTheTheTranslationDialog = true
	}

	title, err := c.GetTitle(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get title: %w", err)
	}

	// reviews, err := c.GetReviews(page)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get reviews: %w", err)
	// }
	roomInfo, err := c.getRoomInfo(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get room info: %w", err)
	}

	photos, err := c.getPhotos(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get photos: %w", err)
	}
	err = page.WaitIdle(defaultWaitTime)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for page to load after getting photos: %w", err)
	}

	// Extract description
	description, err := c.getDescription(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get description: %w", err)
	}

	// Convert photo URLs to strings
	photoStrings := make([]string, len(photos))
	for i, photoURL := range photos {
		photoStrings[i] = photoURL.String()
	}

	return &Listing{
		URL:         parsedURL.String(),
		Title:       title,
		Description: description,
		Photos:      photoStrings,
		RoomInfo:    roomInfo,
	}, nil
}

func (c *Client) getDescription(page *rod.Page) ([]string, error) {
	descButtonSearch, err := page.Timeout(4 * time.Second).Search("div[data-section-id='DESCRIPTION_DEFAULT'] > div > button")
	if err != nil {
		return nil, fmt.Errorf("failed to find description button: %w", err)
	}
	descButton := descButtonSearch.First
	_ = descButton.WaitStable(defaultWaitTime)
	_, err = descButton.CancelTimeout().WaitInteractable()
	if err != nil {
		return nil, fmt.Errorf("failed to wait for description button to be visible: %w", err)
	}
	if err = descButton.CancelTimeout().Timeout(defaultWaitTime).Click("left", 1); err != nil {
		return nil, fmt.Errorf("failed to click description button: %w", err)
	}

	descriptionModal, err := page.Timeout(defaultWaitTime).Element("div[data-section-id='DESCRIPTION_MODAL']")
	if err != nil {
		return nil, fmt.Errorf("failed to find description modal: %w", err)
	}

	sections, err := descriptionModal.CancelTimeout().Elements("section")
	if err != nil {
		return nil, fmt.Errorf("failed to find description sections: %w", err)
	}

	var description []string

	for _, section := range sections {
		// Check if this section has a header with "Registration Details"
		// If so, we've reached the end of the description
		h2, err2 := section.Element("h2")
		if err2 == nil {
			headerText, err3 := h2.Text()
			if err3 == nil && strings.Contains(headerText, "Registration Details") {
				break
			}
		}

		span, err2 := section.Element("span")
		if err2 != nil {
			continue // Skip sections without spans
		}

		text, err2 := span.Text()
		if err2 != nil {
			continue
		}

		if text != "" {
			description = append(description, text)
		}
	}

	return description, nil
}
