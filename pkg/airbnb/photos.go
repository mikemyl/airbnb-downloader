package airbnb

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

const (
	defaultWaitTime = 5 * time.Second
	shortWaitTime   = 500 * time.Millisecond
)

func (c *Client) getPhotos(page *rod.Page) ([]*url.URL, error) {
	err := showAllPhotos(page)
	if err != nil {
		return nil, err
	}

	err2 := clickFirstPhoto(page)
	if err2 != nil {
		return nil, err2
	}

	slideShowDiv, err := page.Timeout(defaultWaitTime).Element("div[data-testid='photo-viewer-slideshow-desktop']")
	if err != nil {
		return nil, fmt.Errorf("failed to find slideshow div: %w", err)
	}
	if err = slideShowDiv.WaitLoad(); err != nil {
		return nil, fmt.Errorf("failed to wait for slideshow to load: %w", err)
	}

	photos, err := parsePhotoUrls(page)
	if err != nil {
		return nil, fmt.Errorf("failed to parse photo urls: %w", err)
	}

	err = closeAndGoBackToMainPage(page)
	if err != nil {
		return nil, fmt.Errorf("failed to close slideshow: %w", err)
	}

	return photos, nil
}

func showAllPhotos(page *rod.Page) error {
	div, err := page.Timeout(defaultWaitTime).ElementR("button > div > div", "Show all photos")
	if err != nil {
		return fmt.Errorf("failed to find 'Show all photos' button: %w", err)
	}

	parent, err := div.CancelTimeout().Parent()
	if err != nil {
		return fmt.Errorf("failed to find parent of 'Show all photos' button: %w", err)
	}
	button, err := parent.Parent()
	if err != nil {
		return fmt.Errorf("failed to find parent/parent of 'Show all photos' button: %w", err)
	}
	if err = button.Timeout(2*defaultWaitTime).Click("left", 1); err != nil {
		return fmt.Errorf("failed to click 'Show all photos' button: %w", err)
	}
	return nil
}

func closeTranslationOnDialog(page *rod.Page) error {
	modal, err := page.Timeout(2 * defaultWaitTime).Search("div[data-testid='modal-container']")
	if err != nil {
		return fmt.Errorf("failed to find close translation modalContainer: %w", err)
	}
	modal.First.MustType(input.Escape)
	time.Sleep(1 * time.Second)
	err = page.Timeout(4 * time.Second).WaitLoad()
	if err != nil {
		return fmt.Errorf("failed to wait for page to load after closing translation modal: %w", err)
	}
	return nil
}

func clickFirstPhoto(page *rod.Page) error {
	photoViewerSection, err := page.Timeout(defaultWaitTime).Element("div[data-testid='photo-viewer-section']")
	if err != nil {
		return fmt.Errorf("failed to find photo viewer section: %w", err)
	}

	firstImageButton, err := photoViewerSection.CancelTimeout().Timeout(defaultWaitTime).Element("div > button")
	if err != nil {
		return fmt.Errorf("failed to find first image button: %w", err)
	}
	if err = firstImageButton.CancelTimeout().Timeout(defaultWaitTime).Click("left", 1); err != nil {
		return fmt.Errorf("failed to click first image button: %w", err)
	}
	return nil
}

func closeAndGoBackToMainPage(page *rod.Page) error {
	closeButton, err := page.Timeout(defaultWaitTime).ElementR("button", "Close")
	if err != nil {
		return fmt.Errorf("failed to find close button: %w", err)
	}
	if err = closeButton.CancelTimeout().Click("left", 1); err != nil {
		return fmt.Errorf("failed to click close button: %w", err)
	}

	if err = page.NavigateBack(); err != nil {
		return fmt.Errorf("failed to navigate back: %w", err)
	}
	if err = page.Timeout(defaultWaitTime).WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for page to load after navigating back: %w", err)
	}
	return nil
}

func parsePhotoUrls(page *rod.Page) ([]*url.URL, error) {
	var photos []*url.URL

	// Iterate through all images in the slideshow
	hasMorePhotos := true
	for hasMorePhotos {
		imageElementSearch, err := page.Timeout(defaultWaitTime).Search("div[data-testid='photo-viewer-slideshow-desktop'] img")
		if err != nil {
			return nil, fmt.Errorf("failed to find image element: %w", err)
		}
		imageElement := imageElementSearch.First
		imageSrc, err := imageElement.CancelTimeout().Attribute("src")
		if err != nil {
			return nil, fmt.Errorf("failed to get image src: %w", err)
		}
		if imageSrc == nil {
			return nil, errors.New("image src is nil")
		}

		sourceWithoutQueryParams := strings.Split(*imageSrc, "?")[0]

		photoURL, err := url.Parse(sourceWithoutQueryParams)
		if err != nil {
			return nil, fmt.Errorf("failed to parse photo url: %w", err)
		}

		photos = append(photos, photoURL)

		buttonSearch, err := page.Timeout(shortWaitTime * 2).Search("div[data-testid='modal-container'] button[aria-label='Next']")
		if err != nil {
			hasMorePhotos = false
			continue
		}

		button := buttonSearch.First

		if err = button.CancelTimeout().Click("left", 1); err != nil {
			return nil, fmt.Errorf("failed to click next button: %w", err)
		}
	}
	return photos, nil
}
