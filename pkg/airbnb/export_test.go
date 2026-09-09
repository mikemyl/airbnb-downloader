package airbnb

import "net/url"

// PhotoURLsFromPageState exposes photoURLsFromPageState to the airbnb_test package.
func PhotoURLsFromPageState(raw []byte) ([]*url.URL, error) {
	return photoURLsFromPageState(raw)
}

// ParseGuestFavoriteReviewHeader exposes parseGuestFavoriteReviewHeader's figures to the airbnb_test package.
func ParseGuestFavoriteReviewHeader(text string, locale Locale) (float64, int, error) {
	summary, err := parseGuestFavoriteReviewHeader(text, locale)
	return summary.score, summary.numberOfReviews, err
}

// ParseReviewHeader exposes parseReviewHeader's figures to the airbnb_test package.
func ParseReviewHeader(text string, locale Locale) (float64, int, bool, error) {
	summary, err := parseReviewHeader(text, locale)
	return summary.score, summary.numberOfReviews, summary.rated, err
}
