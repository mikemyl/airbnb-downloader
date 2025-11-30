## Airbnb Downloader

This project aims to download the basic information of our AirBnb listings and store it locally.

[![Go Report Card](https://goreportcard.com/badge/github.com/mikemyl/airbnb-downloader)](https://goreportcard.com/report/github.com/mikemyl/airbnb-downloader)

AirBnB does not provide an API for this, and getting that info for each one of our listings manually, is a pain. I therefore created this tool to automate the process.

Currently, it extracts:
- Title
- Description (split into paragraphs)
- Photos
- Reviews
- RoomInfo
  * Number of guests
  * Number of bedrooms
  * Number of beds
  * Number of baths

### Dependencies

- [Rod](https://github.com/go-rod/rod).

### How to run

Download the latest release from the [releases](https://github.com/mikemyl/airbnb-downloader/releases) page, for your platform.

Example usage:
```bash
airbnb-downloader https://www.airbnb.com/rooms/12345
```

Also, check the [examples](examples) directory for integrating it as a library.


