package airbnb

import "net/url"

type RoomInfo struct {
	NumberOfGuests   int
	NumberOfBedrooms int
	NumberOfBeds     int
	NumberOfBaths    int
}

// Listing represents an Airbnb listing with its metadata.
type Listing struct {
	URL         *url.URL
	Title       string
	RoomInfo    *RoomInfo
	Description []string
	Photos      []*url.URL
}
