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
		"https://www.airbnb.com/rooms/49703920?guests=1&adults=1&s=67&unique_share_id=2c54d9f1-2a65-43e7-adf9-8192433d443f",
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
		log.Printf("Bedrooms: %d\n", listing.RoomInfo.NumberOfBedrooms)
		log.Printf("Beds: %d\n", listing.RoomInfo.NumberOfBeds)
		log.Printf("Baths: %d\n", listing.RoomInfo.NumberOfBaths)
	}
}
