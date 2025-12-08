package airbnb

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

func (c *Client) getDescription(page *rod.Page) ([]string, error) {
	descButtonSearch, err := page.Timeout(4 * time.Second).Search("div[data-section-id='DESCRIPTION_DEFAULT'] > div > button")
	if err != nil {
		return maybeReadDescriptionSpan(page)
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

	err = closeDescriptionAndGoBackToMainPage(page)
	if err != nil {
		return nil, fmt.Errorf("failed to close description and go back to main page: %w", err)
	}

	return description, nil
}

func maybeReadDescriptionSpan(page *rod.Page) ([]string, error) {
	descriptionSpanSearch, err := page.Timeout(4 * time.Second).Search("div[data-section-id='DESCRIPTION_DEFAULT'] span")
	if err != nil {
		return nil, fmt.Errorf("failed to find description button or span: %w", err)
	}
	descriptionSpan := descriptionSpanSearch.First
	text, err := descriptionSpan.Text()
	if err != nil {
		return nil, fmt.Errorf("failed to get description span text: %w", err)
	}
	textWithoutRegistrationDetails := strings.Split(text, "Registration Details")[0]
	return []string{strings.TrimSpace(textWithoutRegistrationDetails)}, nil
}

func closeDescriptionAndGoBackToMainPage(page *rod.Page) error {
	err := page.NavigateBack()
	if err != nil {
		return fmt.Errorf("failed to navigate back: %w", err)
	}
	if err = page.Timeout(defaultWaitTime).WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for page to load after navigating back: %w", err)
	}
	return nil
}
