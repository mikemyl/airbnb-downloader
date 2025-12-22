package airbnb

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
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
	if _, err = slideShowDiv.WaitInteractable(); err != nil {
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
	searchResult, err := page.Timeout(defaultWaitTime).Search("div[data-section-id='HERO_DEFAULT'] button")
	if err != nil {
		return fmt.Errorf("failed to find 'Show all photos' button: %w", err)
	}
	allButtons, err := searchResult.All()
	if err != nil {
		return fmt.Errorf("failed to get all buttons: %w", err)
	}

	err = allButtons.Last().CancelTimeout().Timeout(defaultWaitTime).Click("left", 1)
	if err != nil {
		return fmt.Errorf("failed to click 'Show all photos' button: %w", err)
	}
	return nil
}

func closeTranslationOnDialog(page *rod.Page) error {
	modal, err := page.Timeout(2 * defaultWaitTime).Search("div[data-testid='modal-container'] div[role='dialog'] button")
	if err != nil {
		return fmt.Errorf("failed to find close translation modalContainer: %w", err)
	}
	err = modal.First.CancelTimeout().Timeout(shortWaitTime).Click("left", 1)
	if err != nil {
		return fmt.Errorf("failed to click close translation modal: %w", err)
	}
	buttons, err := page.Timeout(defaultWaitTime).Search("div[data-testid='main-cookies-banner-container'] button")
	if err != nil {
		return fmt.Errorf("failed to find close cookies banner button: %w", err)
	}
	allButtons, err := buttons.All()
	if err != nil {
		return fmt.Errorf("failed to get all buttons: %w", err)
	}
	err = allButtons.Last().CancelTimeout().Timeout(defaultWaitTime).Click("left", 1)
	if err != nil {
		return fmt.Errorf("failed to click close cookies banner: %w", err)
	}
	err = page.Timeout(2 * time.Second).WaitLoad()
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
	searchResults, err := page.Timeout(defaultWaitTime).Search("div[data-testid='modal-container'] button > span > span > svg")
	if err != nil {
		return fmt.Errorf("failed to find close button: %w", err)
	}
	svg := searchResults.First
	firstSpan, err := svg.CancelTimeout().Parent()
	if err != nil {
		return fmt.Errorf("failed to get parent of svg: %w", err)
	}
	secondSpan, err := firstSpan.Parent()
	if err != nil {
		return fmt.Errorf("failed to get parent of first span: %w", err)
	}
	closeButton, err := secondSpan.Parent()
	if err != nil {
		return fmt.Errorf("failed to get parent of second span: %w", err)
	}
	if err = closeButton.Timeout(defaultWaitTime).Click("left", 1); err != nil {
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

	hasMorePhotos := true
	for hasMorePhotos {
		_ = page.WaitIdle(defaultWaitTime)
		imageElementSearch, err := page.Timeout(defaultWaitTime * 2).Search("div[data-testid='photo-viewer-slideshow-desktop'] img")
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

		buttonSearch, err := page.Timeout(defaultWaitTime).Search("div[data-testid='photo-viewer-slideshow-desktop'] button")
		if err != nil {
			hasMorePhotos = false
			continue
		}

		allButtons, err := buttonSearch.All()
		if err != nil || (len(photos) > 1 && buttonSearch.ResultCount == 1) {
			hasMorePhotos = false
			continue
		}
		nextButton := allButtons.Last().CancelTimeout().Timeout(shortWaitTime)
		if nextButton == nil {
			hasMorePhotos = false
			continue
		}

		if err = nextButton.CancelTimeout().Click("left", 1); err != nil {
			return nil, fmt.Errorf("failed to click next button: %w", err)
		}
	}
	return photos, nil
}
