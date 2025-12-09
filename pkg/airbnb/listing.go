package airbnb

import (
	"fmt"
	"net/url"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// GetListing extracts all information from an Airbnb listing page.
func (c *Client) GetListing(listingURL string) (*Listing, error) {
	parsedURL, err := url.Parse(listingURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse listing url: %w", err)
	}

	target := proto.TargetCreateTarget{
		URL:                     listingURL,
		Width:                   nil,
		Height:                  nil,
		BrowserContextID:        "",
		EnableBeginFrameControl: false,
		NewWindow:               false,
		Background:              false,
		ForTab:                  false,
	}
	page, err := c.browser.Page(target)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}
	defer func(page *rod.Page) {
		_ = page.Close()
	}(page)

	if err = page.WaitLoad(); err != nil {
		return nil, fmt.Errorf("failed to wait for page to load: %w", err)
	}

	if !c.hasGonePastTheTheTranslationDialog {
		_ = closeTranslationOnDialog(page)
		c.hasGonePastTheTheTranslationDialog = true
	}

	title, err := c.GetTitle(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get title: %w", err)
	}

	// reviews, err := c.getReviews(page)
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
	photos = removeDuplicates(photos)
	err = page.WaitIdle(defaultWaitTime)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for page to load after getting photos: %w", err)
	}

	// Extract description
	description, err := c.getDescription(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get description: %w", err)
	}

	amenities, err := c.getAmenities(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get amenities: %w", err)
	}

	reviews, err := c.getReviews(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
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
		Amenities:   amenities,
		Reviews:     reviews,
	}, nil
}

func removeDuplicates(photos []*url.URL) []*url.URL {
	noDups := make([]*url.URL, 0, len(photos))
	seen := make(map[string]bool)
	for _, photo := range photos {
		if !seen[photo.String()] {
			seen[photo.String()] = true
			noDups = append(noDups, photo)
		}
	}
	return noDups
}
