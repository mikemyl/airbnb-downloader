package airbnb

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func (c *Client) GetAmenities(listingURL string) ([]string, error) {
	_, err := url.Parse(listingURL)
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

	amenities, err := c.getAmenities(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get amenities: %w", err)
	}

	return amenities, nil
}

func (c *Client) getAmenities(page *rod.Page) ([]string, error) {
	amenitiesModal, err := openAmenitiesModal(page)
	if err != nil {
		return nil, err
	}

	amenitiesDivs, err := amenitiesModal.CancelTimeout().Elements("section > section > div")
	if err != nil {
		return nil, fmt.Errorf("failed to find amenities amenitiesDivs: %w", err)
	}

	var amenities []string

	for _, div := range amenitiesDivs {
		// Check if this div has a header with "Registration Details"
		// If so, we've reached the end of the amenities
		amenityCategoryEl, err2 := div.Element("h2")
		if err2 != nil {
			return nil, fmt.Errorf("failed to find amenityCategoryEl in amenities div: %w", err2)
		}
		amenityCategory, err2 := amenityCategoryEl.Text()
		if err2 != nil {
			return nil, fmt.Errorf("failed to get amenityCategory text: %w", err2)
		}
		if amenityCategory == "Not included" {
			break
		}

		amenityListItems, err2 := div.Elements("li")
		if err2 != nil {
			return nil, fmt.Errorf("failed to find amenity list items: %w", err2)
		}

		for _, item := range amenityListItems {
			textDiv, err3 := item.Elements("div > div > div")
			if err3 != nil {
				return nil, fmt.Errorf("failed to find text div in amenity list item: %w", err3)
			}
			lastDiv := textDiv[len(textDiv)-1]
			text, err3 := lastDiv.Text()
			if err3 != nil {
				return nil, fmt.Errorf("failed to get text from text div in amenity list item: %w", err3)
			}
			amenities = append(amenities, text)
		}
	}
	err = navigateBack(page)
	if err != nil {
		return nil, fmt.Errorf("failed to close amenities and go back to main page: %w", err)
	}

	return amenities, nil
}

func openAmenitiesModal(page *rod.Page) (*rod.Element, error) {
	amenitiesButtonSearch, err := page.Timeout(4 * time.Second).Search("div[data-section-id='AMENITIES_DEFAULT'] div > button")
	if err != nil {
		return nil, fmt.Errorf("failed to find amenities button: %w", err)
	}
	amenitiesButton := amenitiesButtonSearch.First
	_ = amenitiesButton.WaitStable(defaultWaitTime)
	_, err = amenitiesButton.CancelTimeout().WaitInteractable()
	if err != nil {
		return nil, fmt.Errorf("failed to wait for amenities button to be visible: %w", err)
	}
	if err = amenitiesButton.CancelTimeout().Timeout(defaultWaitTime).Click("left", 1); err != nil {
		return nil, fmt.Errorf("failed to click amenities button: %w", err)
	}

	amenitiesModal, err := page.Timeout(defaultWaitTime).Element("div[data-testid='modal-container'] section")
	if err != nil {
		return nil, fmt.Errorf("failed to find amenities secion in modal: %w", err)
	}
	return amenitiesModal, nil
}
