package airbnb

type RoomInfo struct {
	NumberOfGuests   int `json:"numberOfGuests"`
	NumberOfBedrooms int `json:"numberOfBedrooms"`
	NumberOfBeds     int `json:"numberOfBeds"`
	NumberOfBaths    int `json:"numberOfBaths"`
}

// Listing represents an Airbnb listing with its metadata.
type Listing struct {
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	RoomInfo    *RoomInfo `json:"roomInfo"`
	Description []string  `json:"description"`
	Photos      []string  `json:"photos"`
}
