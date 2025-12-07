package airbnb

type RoomInfo struct {
	NumberOfGuests   int     `json:"numberOfGuests"`
	NumberOfBedrooms float64 `json:"numberOfBedrooms"`
	NumberOfBeds     int     `json:"numberOfBeds"`
	NumberOfBaths    float64 `json:"numberOfBaths"`
}

// Listing represents an Airbnb listing with its metadata.
type Listing struct {
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	RoomInfo    *RoomInfo `json:"roomInfo"`
	Description []string  `json:"description"`
	Photos      []string  `json:"photos"`
}
