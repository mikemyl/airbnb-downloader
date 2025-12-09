package airbnb

import (
	"fmt"

	"github.com/go-rod/rod"
)

func navigateBack(page *rod.Page) error {
	err := page.NavigateBack()
	if err != nil {
		return fmt.Errorf("failed to navigate back: %w", err)
	}
	if err = page.Timeout(defaultWaitTime).WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for page to load after navigating back: %w", err)
	}
	return nil
}
