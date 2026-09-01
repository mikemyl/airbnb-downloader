package airbnb_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mikemyl/airbnb-downloader/pkg/airbnb"
)

const ratedReviewsSection = `
<div data-section-id="REVIEWS_DEFAULT">
  <h2><div><span>4.85 &middot; 146 reviews</span></div></h2>
  <div data-testid="content-scroller">
    <div><div>Cleanliness</div><div>4.9</div></div>
    <div><div>Accuracy</div><div>4.8</div></div>
    <div><div>Communication</div><div>4.7</div></div>
    <div><div>Location</div><div>4.6</div></div>
    <div><div>Check-in</div><div>4.5</div></div>
    <div><div>Value</div><div>4.4</div></div>
  </div>
</div>`

const guestFavoriteReviewsSection = `
<div data-section-id="REVIEWS_DEFAULT">
  <div><h2><span>Rated 4.85 out of 5 from 1,234 reviews.</span></h2></div>
  <div data-testid="content-scroller">
    <div><div>Cleanliness</div><div>4.9</div></div>
    <div><div>Accuracy</div><div>4.8</div></div>
    <div><div>Communication</div><div>4.7</div></div>
    <div><div>Location</div><div>4.6</div></div>
    <div><div>Check-in</div><div>4.5</div></div>
    <div><div>Value</div><div>4.4</div></div>
  </div>
</div>`

const unratedReviewsSection = `
<div data-section-id="REVIEWS_DEFAULT">
  <h2><div><span>New &middot; 2 reviews</span></div></h2>
</div>`

const singleReviewSection = `
<div data-section-id="REVIEWS_DEFAULT">
  <h2><div><span>New &middot; 1 review</span></div></h2>
</div>`

const emptyReviewsSection = `
<div data-section-id="REVIEWS_EMPTY_DEFAULT">
  <div>New &middot; No reviews (yet)</div>
  <div>Be the first to review this place.</div>
</div>`

const greekRatedReviewsSection = `
<div data-section-id="REVIEWS_DEFAULT">
  <h2><div><span>4,85 &middot; 146 κριτικές</span></div></h2>
  <div data-testid="content-scroller">
    <div><div>Καθαριότητα</div><div>4,9</div></div>
    <div><div>Ακρίβεια</div><div>4,8</div></div>
    <div><div>Επικοινωνία</div><div>4,7</div></div>
    <div><div>Τοποθεσία</div><div>4,6</div></div>
    <div><div>Άφιξη</div><div>4,5</div></div>
    <div><div>Τιμή</div><div>4,4</div></div>
  </div>
</div>`

const greekUnratedReviewsSection = `
<div data-section-id="REVIEWS_DEFAULT">
  <h2><div><span>Νέο &middot; 1 κριτική</span></div></h2>
</div>`

const unparseableReviewCountSection = `
<div data-section-id="REVIEWS_DEFAULT">
  <h2><div><span>4.85 &middot; loads of reviews</span></div></h2>
</div>`

const noReviewsSectionAtAll = `<div data-section-id="AMENITIES_DEFAULT">Amenities</div>`

// One client serves every case: launching Chrome and dismissing the cookie
// banner costs more than the scrapes themselves.
func TestGetReviews(t *testing.T) {
	client, err := airbnb.NewClient(airbnb.WithHeadless(true))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	t.Cleanup(func() {
		_ = client.Close()
	})

	fullScores := airbnb.Reviews{
		Rated: true, Score: 4.85, NumberOfReviews: 146,
		ScoreCleanliness: 4.9, ScoreAccuracy: 4.8, ScoreCommunication: 4.7,
		ScoreLocation: 4.6, ScoreCheckIn: 4.5, ScoreValue: 4.4,
	}
	guestFavoriteScores := fullScores
	guestFavoriteScores.NumberOfReviews = 1234

	tests := []struct {
		name        string
		section     string
		locale      airbnb.Locale
		want        airbnb.Reviews
		wantErrLike string
	}{
		{
			name:    "listing with an average rating",
			section: ratedReviewsSection,
			locale:  airbnb.English,
			want:    fullScores,
		},
		{
			name:    "guest favourite listing, whose header names the count in thousands",
			section: guestFavoriteReviewsSection,
			locale:  airbnb.English,
			want:    guestFavoriteScores,
		},
		{
			name:    "listing with too few reviews to be rated",
			section: unratedReviewsSection,
			locale:  airbnb.English,
			want:    airbnb.Reviews{NumberOfReviews: 2},
		},
		{
			name:    "listing with a single review",
			section: singleReviewSection,
			locale:  airbnb.English,
			want:    airbnb.Reviews{NumberOfReviews: 1},
		},
		{
			name:    "listing with no reviews at all",
			section: emptyReviewsSection,
			locale:  airbnb.English,
			want:    airbnb.Reviews{},
		},
		{
			name:    "Greek listing with an average rating",
			section: greekRatedReviewsSection,
			locale:  airbnb.Greek,
			want:    fullScores,
		},
		{
			name:    "Greek listing with too few reviews to be rated",
			section: greekUnratedReviewsSection,
			locale:  airbnb.Greek,
			want:    airbnb.Reviews{NumberOfReviews: 1},
		},
		{
			name:        "header whose review count is not a number",
			section:     unparseableReviewCountSection,
			locale:      airbnb.English,
			wantErrLike: "failed to parse number of reviews",
		},
		{
			name:        "listing that renders no reviews section this library knows",
			section:     noReviewsSectionAtAll,
			locale:      airbnb.English,
			wantErrLike: "failed to find the reviews section",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, getErr := client.GetLocalizedReviews(listingServer(t, test.section).URL, test.locale)

			if test.wantErrLike != "" {
				assertErrorMentions(t, getErr, test.wantErrLike, got)
				return
			}
			if getErr != nil {
				t.Fatalf("GetLocalizedReviews returned an error: %v", getErr)
			}
			if *got != test.want {
				t.Errorf("GetLocalizedReviews\n got: %+v\nwant: %+v", *got, test.want)
			}
		})
	}
}

func assertErrorMentions(t *testing.T, err error, want string, got *airbnb.Reviews) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error containing %q, got reviews %+v", want, *got)
	}
	if !strings.Contains(err.Error(), want) {
		t.Errorf("error %q does not mention %q", err, want)
	}
}

func listingServer(t *testing.T, reviewsSection string) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "<!doctype html><html><head><title>listing</title></head><body>%s</body></html>", reviewsSection)
	}))
	t.Cleanup(server.Close)
	return server
}
