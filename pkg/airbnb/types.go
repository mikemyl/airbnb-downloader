package airbnb

import "net/url"

// Listing represents an Airbnb listing with its metadata.
type Listing struct {
	URL         *url.URL
	Description []string
	Photos      []*url.URL
}
