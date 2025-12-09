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
	_, err := url.Parse(listingURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse listing url: %w", err)
	}

	target := proto.TargetCreateTarget{URL: listingURL}
	page, err := c.browser.Page(target)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}
	defer func(page *rod.Page) {
		_ = page.Close()
	}(page)

	if err = page.WaitLoad(); err != nil {
		return nil, err
	}

	if !c.hasGonePastTheTheTranslationDialog {
		_ = closeTranslationOnDialog(page)
		c.hasGonePastTheTheTranslationDialog = true
	}

	reviews, err := c.getReviews(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}

	return reviews, nil
}

func (c *Client) getReviews(page *rod.Page) (*Reviews, error) {
	searchResults, err := page.Timeout(defaultWaitTime).Search("div[data-section-id='REVIEWS_DEFAULT']")
	if err != nil {
		return nil, fmt.Errorf("failed to find reviews link: %w", err)
	}
	reviewsDiv := searchResults.First
	scoreAndNumberOfReviews, err := reviewsDiv.CancelTimeout().Timeout(defaultWaitTime).Element("h2 > div > span")
	if err != nil {
		return nil, fmt.Errorf("failed to find score and number of reviews: %w", err)
	}
	scoreAndNumberOfReviewsText, err := scoreAndNumberOfReviews.Text()
	if err != nil {
		return nil, fmt.Errorf("failed to get score and number of reviews text: %w", err)
	}
	removedReviewsText := strings.Split(scoreAndNumberOfReviewsText, " reviews")[0]
	scoreAndNumberOfReviewsParts := strings.Split(removedReviewsText, " · ")
	if len(scoreAndNumberOfReviewsParts) != 2 {
		return nil, fmt.Errorf("failed to parse score and number of reviews: %s", scoreAndNumberOfReviewsText)
	}
	var reviews Reviews
	reviews.Score, err = strconv.ParseFloat(scoreAndNumberOfReviewsParts[0], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse score: %w", err)
	}
	reviews.NumberOfReviews, err = strconv.Atoi(scoreAndNumberOfReviewsParts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse number of reviews: %w", err)
	}

	searchResults, err = page.Timeout(defaultWaitTime).Search("div[data-section-id='REVIEWS_DEFAULT'] div[data-testid='content-scroller']")
	if err != nil {
		return nil, fmt.Errorf("failed to find reviews scroller: %w", err)
	}
	scroller := searchResults.First.CancelTimeout()
	reviews.ScoreCleanliness, err = getReviewScore("Cleanliness", scroller)
	if err != nil {
		return nil, fmt.Errorf("failed to get cleanliness review score: %w", err)
	}
	reviews.ScoreAccuracy, err = getReviewScore("Accuracy", scroller)
	if err != nil {
		return nil, fmt.Errorf("failed to get accuracy review score: %w", err)
	}
	reviews.ScoreCommunication, err = getReviewScore("Communication", scroller)
	if err != nil {
		return nil, fmt.Errorf("failed to get communication review score: %w", err)
	}
	reviews.ScoreLocation, err = getReviewScore("Location", scroller)
	if err != nil {
		return nil, fmt.Errorf("failed to get location review score: %w", err)
	}
	reviews.ScoreCheckIn, err = getReviewScore("Check-in", scroller)
	if err != nil {
		return nil, fmt.Errorf("failed to get check-in review score: %w", err)
	}
	reviews.ScoreValue, err = getReviewScore("Value", scroller)
	if err != nil {
		return nil, fmt.Errorf("failed to get value review score: %w", err)
	}
	return &reviews, nil
}

func getReviewScore(reviewType string, element *rod.Element) (float64, error) {
	reviewElement, err := element.ElementR("div", reviewType)
	if err != nil {
		return 0, fmt.Errorf("failed to find %s review: %w", reviewType, err)
	}
	parent, err := reviewElement.Parent()
	if err != nil {
		return 0, fmt.Errorf("failed to find parent of %s review: %w", reviewType, err)
	}
	children, err := parent.Elements("div")
	if err != nil {
		return 0, fmt.Errorf("failed to find children of parent of %s review: %w", reviewType, err)
	}
	lastChild := children[len(children)-1]
	text, err := lastChild.Text()
	if err != nil {
		return 0, fmt.Errorf("failed to get text of last child of parent of %s review: %w", reviewType, err)
	}
	score, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s review score: %w", reviewType, err)
	}
	return score, nil
}
