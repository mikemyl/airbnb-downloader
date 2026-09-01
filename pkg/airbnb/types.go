package airbnb

type RoomInfo struct {
	NumberOfGuests   int     `json:"numberOfGuests"`
	NumberOfBedrooms float64 `json:"numberOfBedrooms"`
	NumberOfBeds     int     `json:"numberOfBeds"`
	NumberOfBaths    float64 `json:"numberOfBaths"`
}

type Reviews struct {
	// Rated is false while Airbnb withholds the average rating, which it does
	// until a listing has three reviews. Score and the per-category scores are
	// then meaningless and NumberOfReviews is the only figure Airbnb publishes.
	Rated              bool    `json:"rated"`
	Score              float64 `json:"scoreTotal"`
	NumberOfReviews    int     `json:"numberOfReviews"`
	ScoreCleanliness   float64 `json:"scoreCleanliness"`
	ScoreAccuracy      float64 `json:"scoreAccuracy"`
	ScoreCommunication float64 `json:"scoreCommunication"`
	ScoreLocation      float64 `json:"scoreLocation"`
	ScoreCheckIn       float64 `json:"scoreCheckIn"`
	ScoreValue         float64 `json:"scoreValue"`
}

// Listing represents an Airbnb listing with its metadata.
type Listing struct {
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	RoomInfo    *RoomInfo `json:"roomInfo"`
	Description []string  `json:"description"`
	Photos      []string  `json:"photos"`
	Amenities   []string  `json:"amenities"`
	Reviews     *Reviews  `json:"reviews"`
}

type Locale string

const (
	English Locale = "en"
	Greek   Locale = "el"
)
