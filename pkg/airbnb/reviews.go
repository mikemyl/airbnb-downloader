package airbnb

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// scoreSeparator divides the score from the review count in the reviews header.
const scoreSeparator = " · "

func (c *Client) GetReviews(listingURL string) (*Reviews, error) {
	return c.GetLocalizedReviews(listingURL, English)
}

// GetLocalizedReviews reads the reviews of a listing rendered in the given locale.
func (c *Client) GetLocalizedReviews(listingURL string, locale Locale) (*Reviews, error) {
	target := proto.TargetCreateTarget{
		URL:                     listingURL,
		Width:                   nil,
		Height:                  nil,
		BrowserContextID:        "",
		EnableBeginFrameControl: false,
		NewWindow:               false,
		Background:              false,
		ForTab:                  false,
	}
	page, err := c.browser.Page(target)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}
	defer func(page *rod.Page) {
		_ = page.Close()
	}(page)

	_, err = url.Parse(listingURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse listing url: %w", err)
	}

	if err = page.WaitLoad(); err != nil {
		return nil, fmt.Errorf("failed to wait for page to load: %w", err)
	}

	if !c.bannerIsClosedMap[locale] {
		_ = closeTranslationOnDialog(page)
		c.bannerIsClosedMap[locale] = true
	}

	reviews, err := c.getReviews(page, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}

	return reviews, nil
}

func (c *Client) getReviews(page *rod.Page, locale Locale) (*Reviews, error) {
	summary, err := getReviewSummary(page, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to get score and number of reviews: %w", err)
	}

	reviews := Reviews{
		Rated:              summary.rated,
		Score:              summary.score,
		NumberOfReviews:    summary.numberOfReviews,
		ScoreCleanliness:   0,
		ScoreAccuracy:      0,
		ScoreCommunication: 0,
		ScoreLocation:      0,
		ScoreCheckIn:       0,
		ScoreValue:         0,
	}
	if !summary.rated {
		return &reviews, nil
	}

	searchResults, err := page.Timeout(defaultWaitTime).Search("div[data-section-id='REVIEWS_DEFAULT'] div[data-testid='content-scroller']")
	if err != nil {
		return nil, fmt.Errorf("failed to find reviews scroller: %w", err)
	}
	scroller := searchResults.First.CancelTimeout()
	reviews.ScoreCleanliness, err = getReviewScore(getCleaningnessText(locale), scroller, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to get cleanliness review score: %w", err)
	}
	reviews.ScoreAccuracy, err = getReviewScore(getAccuracyText(locale), scroller, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to get accuracy review score: %w", err)
	}
	reviews.ScoreCommunication, err = getReviewScore(getCommunicationText(locale), scroller, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to get communication review score: %w", err)
	}
	reviews.ScoreLocation, err = getReviewScore(getLocationText(locale), scroller, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to get location review score: %w", err)
	}
	reviews.ScoreCheckIn, err = getReviewScore(getCheckInText(locale), scroller, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to get check-in review score: %w", err)
	}
	reviews.ScoreValue, err = getReviewScore(getPriceText(locale), scroller, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to get value review score: %w", err)
	}
	return &reviews, nil
}

// reviewSummary is the aggregate rating Airbnb prints in the header of the
// reviews section. An unrated listing carries a review count but no scores.
type reviewSummary struct {
	score           float64
	numberOfReviews int
	rated           bool
}

// getReviewSummary reads whichever of the three review headers the listing
// renders: the plain one, the guest-favourite one, or the empty-state section
// Airbnb swaps in when nobody has reviewed the listing yet. That third section
// carries no figures, so its presence alone is the answer.
func getReviewSummary(page *rod.Page, locale Locale) (reviewSummary, error) {
	summary := unratedSummary()
	_, err := page.Timeout(defaultWaitTime).Race().Element("div[data-section-id='REVIEWS_DEFAULT'] h2 > div > span").Handle(
		func(e *rod.Element) error {
			parsed, parseErr := parseElement(e, locale, parseReviewHeader)
			summary = parsed
			return parseErr
		},
	).Element("div[data-section-id='REVIEWS_DEFAULT'] div > h2 > span").Handle(
		func(e *rod.Element) error {
			parsed, parseErr := parseElement(e, locale, parseGuestFavoriteReviewHeader)
			summary = parsed
			return parseErr
		},
	).Element("div[data-section-id='REVIEWS_EMPTY_DEFAULT']").Handle(
		func(_ *rod.Element) error {
			return nil
		},
	).Do()
	if err != nil {
		return unratedSummary(), fmt.Errorf("failed to find the reviews section: %w", err)
	}
	return summary, nil
}

func parseElement(elem *rod.Element, locale Locale, parse func(string, Locale) (reviewSummary, error)) (reviewSummary, error) {
	text, err := elem.CancelTimeout().Text()
	if err != nil {
		return unratedSummary(), fmt.Errorf("failed to get score and number of reviews text: %w", err)
	}
	return parse(text, locale)
}

// parseGuestFavoriteReviewHeader reads the "Rated 4.85 out of 5 from 146 reviews."
// header. Airbnb has reworded the Greek form more than once ("Έλαβε … στα 5 σε …",
// then "Βαθμολογήθηκε με … στα 5 από …"), so the match is on the figures around
// the fixed "out of 5" token rather than the surrounding words.
func parseGuestFavoriteReviewHeader(text string, locale Locale) (reviewSummary, error) {
	match := guestFavoriteHeaderPattern(locale).FindStringSubmatch(text)
	if match == nil {
		return unratedSummary(), fmt.Errorf("failed to parse score and number of reviews: %s", text)
	}
	return newReviewSummary(match[1], match[2], locale)
}

func guestFavoriteHeaderPattern(locale Locale) *regexp.Regexp {
	switch locale {
	case Greek:
		return greekGuestFavoriteHeader
	case English:
		return englishGuestFavoriteHeader
	default:
		return englishGuestFavoriteHeader
	}
}

var (
	englishGuestFavoriteHeader = regexp.MustCompile(`(\d+(?:\.\d+)?) out of 5 from ([\d,]+) reviews?`)
	greekGuestFavoriteHeader   = regexp.MustCompile(`(\d+(?:,\d+)?) στα 5 (?:σε|από) ([\d.]+) κριτικ`)
)

// parseReviewHeader reads the "4.85 · 146 reviews" header, which reads
// "New · 1 review" until the listing has enough reviews to be rated.
func parseReviewHeader(text string, locale Locale) (reviewSummary, error) {
	parts := strings.Split(text, scoreSeparator)
	if len(parts) == 1 {
		// Listings with too few reviews to be scored render just "2 reviews".
		return newReviewSummary(getNewListingText(locale), parts[0], locale)
	}
	if len(parts) != 2 {
		return unratedSummary(), fmt.Errorf("failed to parse score and number of reviews: %s", text)
	}
	return newReviewSummary(parts[0], parts[1], locale)
}

func newReviewSummary(scoreText, countText string, locale Locale) (reviewSummary, error) {
	numberOfReviews, err := parseNumberOfReviews(countText)
	if err != nil {
		return unratedSummary(), err
	}
	// Anything other than the exact "New" marker still has to parse as a score,
	// so a changed Airbnb label fails loudly instead of reading as unrated.
	if strings.TrimSpace(scoreText) == getNewListingText(locale) {
		return reviewSummary{score: 0, numberOfReviews: numberOfReviews, rated: false}, nil
	}
	score, err := strconv.ParseFloat(translateDecimal(strings.TrimSpace(scoreText), locale), 64)
	if err != nil {
		return unratedSummary(), fmt.Errorf("failed to parse score: %w", err)
	}
	return reviewSummary{score: score, numberOfReviews: numberOfReviews, rated: true}, nil
}

// parseNumberOfReviews reads the count out of "146 reviews", "1 review" or
// "1 κριτική", ignoring any thousands separator.
func parseNumberOfReviews(text string) (int, error) {
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return 0, fmt.Errorf("failed to parse number of reviews: %s", text)
	}
	numberOfReviews, err := strconv.Atoi(strings.NewReplacer(",", "", ".", "").Replace(fields[0]))
	if err != nil {
		return 0, fmt.Errorf("failed to parse number of reviews from %s: %w", text, err)
	}
	return numberOfReviews, nil
}

func unratedSummary() reviewSummary {
	return reviewSummary{score: 0, numberOfReviews: 0, rated: false}
}

func getReviewScore(reviewType string, element *rod.Element, locale Locale) (float64, error) {
	// Use regex anchors to match exactly the review type text
	pattern := "^" + reviewType + "$"
	reviewElement, err := element.ElementR("div", pattern)
	if err != nil {
		return 0, fmt.Errorf("failed to find %s review: %w", reviewType, err)
	}
	sibling, err := reviewElement.Next()
	if err != nil {
		return 0, fmt.Errorf("failed to get next sibling of %s review: %w", reviewType, err)
	}
	text, err := sibling.Text()
	if err != nil {
		return 0, fmt.Errorf("failed to get text of last child of parent of %s review: %w", reviewType, err)
	}
	score, err := strconv.ParseFloat(translateDecimal(text, locale), 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s review score: %w", reviewType, err)
	}
	return score, nil
}
