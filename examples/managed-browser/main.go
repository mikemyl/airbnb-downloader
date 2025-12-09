package main

import (
	"log"
	"runtime"

	"github.com/mikemyl/airbnb-downloader/pkg/airbnb"
)

func main() {
	// Create client
	client, err := airbnb.NewClient(airbnb.WithRodURL("localhost:7317"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	listingUrls := []string{
		"https://www.airbnb.com/rooms/51573244?source_impression_id=p3_1763552182_P3BLEoXunPK0xs6K",
		"https://www.airbnb.com/rooms/51573207?source_impression_id=p3_1764364702_P3PINb00ZWafPed0",
	}
	for _, listingURL := range listingUrls {
		log.Printf("\n\n === Fetching listing: %s === \n\n", listingURL)

		listing, listingErr := client.GetListing(listingURL)
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
