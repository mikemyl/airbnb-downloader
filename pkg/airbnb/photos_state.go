package airbnb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/go-rod/rod"
)

// Airbnb embeds the full, ordered photo tour in the page state JSON
// (<script id="data-deferred-state-0">) as a PhotoTourModalSection whose
// mediaItems carry each photo's original baseUrl. Reading it avoids driving
// the photo modal, whose DOM has changed more than once.
const (
	pageStateSelector    = "#data-deferred-state-0"
	photoTourSectionType = "PhotoTourModalSection"
)

var errNoPhotoTourInState = errors.New("no photo tour section in page state")

func (c *Client) getPhotosFromPageState(page *rod.Page) ([]*url.URL, error) {
	stateEl, err := page.Timeout(shortWaitTime).Element(pageStateSelector)
	if err != nil {
		return nil, fmt.Errorf("page state script not found: %w", err)
	}
	raw, err := stateEl.CancelTimeout().Text()
	if err != nil {
		return nil, fmt.Errorf("failed to read page state: %w", err)
	}
	return photoURLsFromPageState([]byte(raw))
}

func photoURLsFromPageState(raw []byte) ([]*url.URL, error) {
	var state any
	if err := json.Unmarshal(raw, &state); err != nil {
		return nil, fmt.Errorf("failed to decode page state: %w", err)
	}

	section := findPhotoTourSection(state)
	if section == nil {
		return nil, errNoPhotoTourInState
	}

	items, _ := section["mediaItems"].([]any)
	photos := make([]*url.URL, 0, len(items))
	for _, item := range items {
		media, ok := item.(map[string]any)
		if !ok {
			continue
		}
		baseURL, ok := media["baseUrl"].(string)
		if !ok || baseURL == "" {
			continue
		}
		parsed, err := url.Parse(baseURL)
		if err != nil {
			return nil, fmt.Errorf("invalid photo url %q: %w", baseURL, err)
		}
		photos = append(photos, parsed)
	}

	if len(photos) == 0 {
		return nil, errNoPhotoTourInState
	}
	return photos, nil
}

func findPhotoTourSection(node any) map[string]any {
	switch v := node.(type) {
	case map[string]any:
		if v["__typename"] == photoTourSectionType {
			if _, ok := v["mediaItems"].([]any); ok {
				return v
			}
		}
		for _, child := range v {
			if found := findPhotoTourSection(child); found != nil {
				return found
			}
		}
	case []any:
		for _, child := range v {
			if found := findPhotoTourSection(child); found != nil {
				return found
			}
		}
	}
	return nil
}
