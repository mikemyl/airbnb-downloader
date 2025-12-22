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
