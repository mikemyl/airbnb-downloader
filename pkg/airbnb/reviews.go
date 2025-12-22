package airbnb

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func (c *Client) GetReviews(listingURL string) (*Reviews, error) {
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

	if !c.hasGonePastTheTheTranslationDialog {
		_ = closeTranslationOnDialog(page)
		c.hasGonePastTheTheTranslationDialog = true
	}

	reviews, err := c.getReviews(page, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}

	return reviews, nil
}

func (c *Client) getReviews(page *rod.Page, locale Locale) (*Reviews, error) {
	reviews := Reviews{
		Score:              0,
		NumberOfReviews:    0,
		ScoreCleanliness:   0,
		ScoreAccuracy:      0,
		ScoreCommunication: 0,
		ScoreLocation:      0,
		ScoreCheckIn:       0,
		ScoreValue:         0,
	}
	_, err := page.Timeout(defaultWaitTime).Race().Element("div[data-section-id='REVIEWS_DEFAULT'] h2 > div > span").Handle(
		func(e *rod.Element) error {
			score, nReviews, err2 := getScoreAndNumberOfReviews(e, locale)
			if err2 != nil {
				return err2
			}
			reviews.Score = score
			reviews.NumberOfReviews = nReviews
			return nil
		},
	).Element("div[data-section-id='REVIEWS_DEFAULT'] div > h2 > span").Handle(
		func(e *rod.Element) error {
			score, nReviews, err2 := getScoreAndNumberOfReviewsForGuestFavorite(e, locale)
			if err2 != nil {
				return err2
			}
			reviews.Score = score
			reviews.NumberOfReviews = nReviews
			return nil
		},
	).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get score and number of reviews: %w", err)
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
	reviews.ScoreCheckIn, err = getReviewScore(getCommunicationText(locale), scroller, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to get check-in review score: %w", err)
	}
	reviews.ScoreValue, err = getReviewScore(getPriceText(locale), scroller, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to get value review score: %w", err)
	}
	return &reviews, nil
}

func getScoreAndNumberOfReviewsForGuestFavorite(elem *rod.Element, locale Locale) (float64, int, error) {
	scoreAndNumberOfReviewsText, err := elem.CancelTimeout().Text()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get score and number of reviews text: %w", err)
	}
	removedRatedText := strings.ReplaceAll(scoreAndNumberOfReviewsText, getRatedText(locale), "")
	removedReviewsText := strings.ReplaceAll(removedRatedText, getReviewsText(locale)+".", "")
	scoreAndNumberOfReviewsParts := strings.Split(removedReviewsText, getOutOfFiveText(locale))
	if len(scoreAndNumberOfReviewsParts) != 2 {
		return 0, 0, fmt.Errorf("failed to parse score and number of reviews: %s", scoreAndNumberOfReviewsText)
	}
	normalizedScore := translateDecimal(scoreAndNumberOfReviewsParts[0], locale)
	score, err := strconv.ParseFloat(normalizedScore, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse score: %w", err)
	}
	reviews, err := strconv.Atoi(scoreAndNumberOfReviewsParts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse number of reviews: %w", err)
	}
	return score, reviews, nil
}

func getScoreAndNumberOfReviews(elem *rod.Element, locale Locale) (float64, int, error) {
	scoreAndNumberOfReviewsText, err := elem.CancelTimeout().Text()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get score and number of reviews text: %w", err)
	}
	removedReviewsText := strings.Split(scoreAndNumberOfReviewsText, getReviewsText(locale))[0]
	scoreAndNumberOfReviewsParts := strings.Split(removedReviewsText, " · ")
	if len(scoreAndNumberOfReviewsParts) != 2 {
		return 0, 0, fmt.Errorf("failed to parse score and number of reviews: %s", scoreAndNumberOfReviewsText)
	}
	score, err := strconv.ParseFloat(translateDecimal(scoreAndNumberOfReviewsParts[0], locale), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse score: %w", err)
	}
	reviews, err := strconv.Atoi(scoreAndNumberOfReviewsParts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse number of reviews: %w", err)
	}
	return score, reviews, nil
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
