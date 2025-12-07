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
	}
}
