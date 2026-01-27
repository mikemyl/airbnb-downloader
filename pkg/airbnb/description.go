package airbnb

import (
	"fmt"
	"strings"

	"github.com/go-rod/rod"
)

func (c *Client) getDescription(page *rod.Page, locale Locale) ([]string, error) {
	descButtonSearch, err := page.Timeout(defaultWaitTime).Search("div[data-section-id='DESCRIPTION_DEFAULT'] > div > button")
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

	description, err := getDescription(descriptionModal, locale)
	if err != nil {
		return nil, err
	}

	err = navigateBack(page)
	if err != nil {
		return nil, fmt.Errorf("failed to close description and go back to main page: %w", err)
	}

	return description, nil
}

func getDescription(descriptionModal *rod.Element, locale Locale) ([]string, error) {
	descriptionModal = descriptionModal.CancelTimeout()
	sections, err := descriptionModal.Timeout(shortWaitTime).Elements("section")
	if err != nil {
		return nil, fmt.Errorf("failed to find description sections: %w", err)
	}

	if len(sections) == 0 {
		descriptionSpan, err2 := descriptionModal.Timeout(shortWaitTime).Element("span")
		if err2 != nil {
			return nil, fmt.Errorf("failed to find description span: %w", err2)
		}
		text, err2 := descriptionSpan.Text()
		if err2 != nil {
			return nil, fmt.Errorf("failed to get description span text: %w", err2)
		}
		textWithoutRegistrationDetails := strings.Split(text, getRegistrationDetailsText(locale))[0]
		textWithoutRegistrationDetails = strings.TrimSpace(textWithoutRegistrationDetails)
		return []string{textWithoutRegistrationDetails}, nil
	}

	var description []string

	for _, section := range sections {
		// Check if this section has a header with "Registration Details"
		// If so, we've reached the end of the description
		h2, err2 := section.Element("h2")
		if err2 == nil {
			headerText, err3 := h2.Text()
			if err3 == nil && strings.Contains(headerText, getRegistrationDetailsText(locale)) {
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

func maybeReadDescriptionSpan(page *rod.Page) ([]string, error) {
	descriptionSpanSearch, err := page.Timeout(defaultWaitTime).Search("div[data-section-id='DESCRIPTION_DEFAULT'] span")
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
