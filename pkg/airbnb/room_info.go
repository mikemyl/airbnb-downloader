package airbnb

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
)

func (c *Client) getRoomInfo(page *rod.Page) (*RoomInfo, error) {
	var roomInfo RoomInfo
	searchResults, err := page.Search("div[data-section-id='OVERVIEW_DEFAULT_V2'] li")
	if err != nil {
		return nil, fmt.Errorf("failed to find room info: %w", err)
	}

	listItems, err := searchResults.All()
	if err != nil {
		return nil, fmt.Errorf("failed to get all list items: %w", err)
	}

	for _, item := range listItems {
		text, err2 := item.Timeout(defaultWaitTime).Text()
		if err2 != nil {
			return nil, fmt.Errorf("failed to get list item text: %w", err2)
		}
		err2 = captureRoomInfo(text, &roomInfo)
		if err2 != nil {
			return nil, fmt.Errorf("failed to capture room info: %w", err2)
		}
	}
	return &roomInfo, nil
}

func captureRoomInfo(text string, r *RoomInfo) error {
	removedDots := strings.ReplaceAll(text, "·", "")
	trimmedText := strings.TrimSpace(removedDots)
	parts := strings.Split(trimmedText, " ")
	if len(parts) == 0 {
		return errors.New("no parts found")
	}

	num, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("failed to convert number to int: %w", err)
	}

	switch {
	case strings.Contains(text, "guest"):
		r.NumberOfGuests = num
	case strings.Contains(text, "bedroom"):
		r.NumberOfBedrooms = num
	case strings.Contains(text, "bed"):
		r.NumberOfBeds = num
	case strings.Contains(text, "bath"):
		r.NumberOfBaths = num
	}

	return nil
}
