package main

import (
	"log"
	"runtime"

	"github.com/mikemyl/airbnb-downloader/pkg/airbnb"
)

func main() {
	// Create client
	client, err := airbnb.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	listingUrls := []string{
		"https://www.airbnb.com/rooms/1438479641258683118",
	}
	for _, listingURL := range listingUrls {
		log.Printf("\n\n === Fetching listing: %s === \n\n", listingURL)

		listing, listingErr := client.GetLocalizedListing(listingURL, airbnb.Greek)
		if listingErr != nil {
			log.Printf("Failed to get listing: %v\n", listingErr)
			runtime.Goexit()
		}

		log.Printf("Title: %s\n", listing.Title)

		log.Println("=== Description ===")
		for i, paragraph := range listing.Description {
			log.Printf("%d. %s\n\n", i+1, paragraph)
		}

		log.Printf("=== Photos (%d total) ===\n", len(listing.Photos))
		for i, photoURL := range listing.Photos {
			log.Printf("%d. %s\n", i+1, photoURL)
		}

		log.Println("=== Room Info ===")
		log.Printf("Guests: %d\n", listing.RoomInfo.NumberOfGuests)
		log.Printf("Bedrooms: %f\n", listing.RoomInfo.NumberOfBedrooms)
		log.Printf("Beds: %d\n", listing.RoomInfo.NumberOfBeds)
		log.Printf("Baths: %f\n", listing.RoomInfo.NumberOfBaths)

		log.Println("=== Amenities ===")
		for i, amenity := range listing.Amenities {
			log.Printf("%d. %s\n", i+1, amenity)
		}

		log.Println("=== Reviews ===")
		log.Printf("Score: %f\n", listing.Reviews.Score)
		log.Printf("Number of Reviews: %d\n", listing.Reviews.NumberOfReviews)
		log.Printf("Score Cleanliness: %f\n", listing.Reviews.ScoreCleanliness)
		log.Printf("Score Accuracy: %f\n", listing.Reviews.ScoreAccuracy)
		log.Printf("Score Communication: %f\n", listing.Reviews.ScoreCommunication)
		log.Printf("Score Location: %f\n", listing.Reviews.ScoreLocation)
		log.Printf("Score Check In: %f\n", listing.Reviews.ScoreCheckIn)
		log.Printf("Score Value: %f\n", listing.Reviews.ScoreValue)
	}
}
