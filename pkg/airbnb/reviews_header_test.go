package airbnb_test

import (
	"testing"

	"github.com/mikemyl/airbnb-downloader/pkg/airbnb"
)

func TestParseGuestFavoriteReviewHeader(t *testing.T) {
	cases := []struct {
		name   string
		text   string
		locale airbnb.Locale
		score  float64
		count  int
	}{
		{"english", "Rated 4.85 out of 5 from 146 reviews.", airbnb.English, 4.85, 146},
		{"english single review", "Rated 5 out of 5 from 1 review.", airbnb.English, 5, 1},
		{"greek old wording", "Έλαβε 4,85 στα 5 σε 146 κριτικές.", airbnb.Greek, 4.85, 146},
		{"greek new wording", "Βαθμολογήθηκε με 5,0 στα 5 από 9 κριτικές.", airbnb.Greek, 5, 9},
		{"greek single review", "Βαθμολογήθηκε με 5,0 στα 5 από 1 κριτική.", airbnb.Greek, 5, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			score, count, err := airbnb.ParseGuestFavoriteReviewHeader(tc.text, tc.locale)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if score != tc.score || count != tc.count {
				t.Errorf("got score=%v count=%d, want score=%v count=%d", score, count, tc.score, tc.count)
			}
		})
	}
}

func TestParseGuestFavoriteReviewHeader_RejectsUnknownWording(t *testing.T) {
	if _, _, err := airbnb.ParseGuestFavoriteReviewHeader("Something else entirely", airbnb.Greek); err == nil {
		t.Fatal("expected an error for an unrecognised header")
	}
}

func TestParseReviewHeader(t *testing.T) {
	cases := []struct {
		name   string
		text   string
		locale airbnb.Locale
		score  float64
		count  int
		rated  bool
	}{
		{"scored", "4.85 · 146 reviews", airbnb.English, 4.85, 146, true},
		{"new listing", "New · 1 review", airbnb.English, 0, 1, false},
		{"too few reviews to score", "2 reviews", airbnb.English, 0, 2, false},
		{"greek too few reviews", "2 κριτικές", airbnb.Greek, 0, 2, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			score, count, rated, err := airbnb.ParseReviewHeader(tc.text, tc.locale)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if score != tc.score || count != tc.count || rated != tc.rated {
				t.Errorf("got score=%v count=%d rated=%v, want score=%v count=%d rated=%v", score, count, rated, tc.score, tc.count, tc.rated)
			}
		})
	}
}

func TestParseReviewHeader_RejectsGarbage(t *testing.T) {
	if _, _, _, err := airbnb.ParseReviewHeader("nothing here", airbnb.English); err == nil {
		t.Fatal("expected an error for an unparseable header")
	}
}
