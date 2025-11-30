package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mikemyl/airbnb-downloader/pkg/airbnb"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// Use custom FlagSet to avoid Rod's global flags
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Define flags
	rodURL := fs.String("rod-url", "", "URL of Rod browser service (e.g., localhost:7317)")
	headless := fs.Bool("headless", true, "Run browser in headless mode")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <urls>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Extract information from Airbnb listings.\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  urls   Single URL or comma-separated list of URLs\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://www.airbnb.com/rooms/12345\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -rod-url localhost:7317 https://www.airbnb.com/rooms/12345\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://www.airbnb.com/rooms/12345,https://www.airbnb.com/rooms/67890\n", os.Args[0])
	}

	fs.Parse(os.Args[1:])

	// Expect exactly one argument
	args := fs.Args()
	if len(args) != 1 {
		fs.Usage()
		os.Exit(1)
	}

	// Split on commas to support multiple URLs
	urls := strings.Split(args[0], ",")

	// Trim whitespace from URLs
	for i := range urls {
		urls[i] = strings.TrimSpace(urls[i])
	}

	// Create client with options
	var opts []airbnb.Option
	if *rodURL != "" {
		opts = append(opts, airbnb.WithRodURL(*rodURL))
	}
	opts = append(opts, airbnb.WithHeadless(*headless))

	client, err := airbnb.NewClient(opts...)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	// Fetch listings
	var listings []*airbnb.Listing
	for _, url := range urls {
		log.Printf("Fetching: %s", url)
		listing, err2 := client.GetListing(url)
		if err2 != nil {
			log.Printf("Error fetching %s: %v", url, err2)
			continue
		}
		listings = append(listings, listing)
	}

	// Output as pretty JSON
	output, err := json.MarshalIndent(listings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	_, err = os.Stdout.Write(output)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	_, err = os.Stdout.Write([]byte("\n"))
	if err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	return nil
}
