package airbnb

import "strings"

func translateDecimal(s string, locale Locale) string {
	switch locale {
	case English:
		return s
	case Greek:
		return strings.ReplaceAll(s, ",", ".")
	default:
		return s
	}
}

func getAccuracyText(locale Locale) string {
	switch locale {
	case English:
		return "Accuracy"
	case Greek:
		return "Ακρίβεια"
	default:
		return "Accuracy"
	}
}

func getCommunicationText(locale Locale) string {
	switch locale {
	case English:
		return "Communication"
	case Greek:
		return "Επικοινωνία"
	default:
		return "Communication"
	}
}

// Airbnb's check-in category label. Verified against a live listing on
// 2026-07-27: "Check-in" (en) / "Άφιξη" (el). getReviewScore anchors this with
// ^...$, so it must match the rendered label exactly.
func getCheckInText(locale Locale) string {
	switch locale {
	case English:
		return "Check-in"
	case Greek:
		return "Άφιξη"
	default:
		return "Check-in"
	}
}

func getLocationText(locale Locale) string {
	switch locale {
	case English:
		return "Location"
	case Greek:
		return "Τοποθεσία"
	default:
		return "Location"
	}
}

func getPriceText(locale Locale) string {
	switch locale {
	case English:
		return "Value"
	case Greek:
		return "Τιμή"
	default:
		return "Value"
	}
}

func getCleaningnessText(locale Locale) string {
	switch locale {
	case English:
		return "Cleanliness"
	case Greek:
		return "Καθαριότητα"
	default:
		return "Cleanliness"
	}
}

func getOutOfFiveText(locale Locale) string {
	switch locale {
	case English:
		return " out of 5 from "
	case Greek:
		return " στα 5 σε "
	default:
		return " out of 5 from "
	}
}

func getReviewsText(locale Locale) string {
	switch locale {
	case English:
		return " reviews"
	case Greek:
		return " κριτικές"
	default:
		return " reviews"
	}
}

// Airbnb's stand-in for the average rating of a listing with fewer than three
// reviews. Verified against a live listing on 2026-09-01: "New" (en) / "Νέο" (el).
func getNewListingText(locale Locale) string {
	switch locale {
	case English:
		return "New"
	case Greek:
		return "Νέο"
	default:
		return "New"
	}
}

func getRatedText(locale Locale) string {
	switch locale {
	case English:
		return "Rated "
	case Greek:
		return "Έλαβε "
	default:
		return "Rated "
	}
}

func getAmenityNotAvailableTranslation(locale Locale) string {
	switch locale {
	case English:
		return "Not included"
	case Greek:
		return "Δεν περιλαμβάνονται"
	default:
		return "Not included"
	}
}

func getRegistrationDetailsText(locale Locale) string {
	switch locale {
	case English:
		return "Registration Details"
	case Greek:
		return "Στοιχεία εγγραφής"
	default:
		return "Registration Details"
	}
}
