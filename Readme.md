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
- Amenities
- Reviews

### Dependencies

- [Rod](https://github.com/go-rod/rod).

### How to run

Download the latest release from the [releases](https://github.com/mikemyl/airbnb-downloader/releases) page, for your platform.

Example usage and output:
```bash
./airbnb-downloader https://www.airbnb.com/rooms/51573207?source_impression_id=p3_1764364702_P3PINb00ZWafPed0

2025/12/01 23:02:26 Fetching: https://www.airbnb.com/rooms/51573207?source_impression_id=p3_1764364702_P3PINb00ZWafPed0
[
  {
    "url": "https://www.airbnb.com/rooms/51573207?source_impression_id=p3_1764364702_P3PINb00ZWafPed0",
    "title": "Voda Penthouse \u0026 Terrace - Acropolis View Rooftop",
    "roomInfo": {
      "numberOfGuests": 2,
      "numberOfBedrooms": 1,
      "numberOfBeds": 1,
      "numberOfBaths": 1
    },
    "description": [
      "Welcome to this bright \u0026 stylish 6th floor, 27 sq.m open plan apartment of Voda Luxury Residence. Relax on the rooftop lounge where you will experience a breathtaking view of Acropolis, Lycabettus and the city of Athens. Newly designed with attention to detail, modern amenities \u0026 200 mbps fiber internet, aiming for an unforgettable stay, Voda Luxury Residence is located only 550 m from Larissa Metro Station, making it very easy to navigate around the city. Check our profile to discover 15 flats!",
      "This open plan studio apartment is 27 sq. m and is located on the 6th floor. It consists of a living area with 2 comfortable armchairs and 43'' Smart TV, a dining area with a stand and stools, a Queen Size Bed with a \"stress free\" mattress and a closet, one bathroom with shower, a fully equipped kitchen and a large balcony. The kitchen includes a fridge/freezer, stove, espresso and filter coffee machine, microwave, kettle as well as all the necessary kitchen and cooking utensils for preparing meals. The apartment is equipped with a secure multi lock entrance door, electrical window shutters and Air Conditioning.",
      "The apartment is located only 550 meters from Larissa Metro Station, 700 meters from Metaxourgeio Metro station \u0026 800 meters from Omonoia Metro Station. All the sightseeing sites can be reached within 5 metro stops or 20-40 minutes on foot."
    ],
    "photos": [
      "https://a0.muscache.com/im/pictures/miso/Hosting-51573207/original/b12bbd2a-e046-4ad9-b04a-f64b2c38d8ef.jpeg",
      "https://a0.muscache.com/im/pictures/miso/Hosting-51573207/original/b12bbd2a-e046-4ad9-b04a-f64b2c38d8ef.jpeg",
      "https://a0.muscache.com/im/pictures/miso/Hosting-51573207/original/b12bbd2a-e046-4ad9-b04a-f64b2c38d8ef.jpeg",
      "https://a0.muscache.com/im/pictures/miso/Hosting-51573207/original/f8a1b6cf-7c10-41f8-b610-6659a57d229a.jpeg",
      "https://a0.muscache.com/im/pictures/miso/Hosting-51573207/original/a4cb12a1-9738-4682-8ad6-1145fac63c93.jpeg",
      "https://a0.muscache.com/im/pictures/miso/Hosting-51573207/original/9e0773fa-91ad-49a0-b7d4-d3256aeca995.jpeg",
      "https://a0.muscache.com/im/pictures/miso/Hosting-51573207/original/f1a1fba7-25aa-4501-853b-0619fd985993.jpeg",
      "https://a0.muscache.com/im/pictures/miso/Hosting-51573207/original/a1616744-6680-4df2-a553-41244e8c742c.jpeg",
      "https://a0.muscache.com/im/pictures/miso/Hosting-51573207/original/180bb426-ff26-4a10-b3f5-62ca3573aae7.jpeg",
      "https://a0.muscache.com/im/pictures/miso/Hosting-51573207/original/a5abfcbb-64ff-46e3-ab2f-1ce679d309ed.jpeg",
      "https://a0.muscache.com/im/pictures/hosting/Hosting-U3RheVN1cHBseUxpc3Rpbmc6NTE1NzMyMDc%3D/original/4c5c33b3-95d6-4689-9942-4b881025df74.jpeg",
      "https://a0.muscache.com/im/pictures/hosting/Hosting-U3RheVN1cHBseUxpc3Rpbmc6NTE1NzMyMDc%3D/original/eb3c0630-4d53-4a29-a781-91973b987986.jpeg",
      "https://a0.muscache.com/im/pictures/hosting/Hosting-U3RheVN1cHBseUxpc3Rpbmc6NTE1NzMyMDc%3D/original/61b9a737-e837-444b-bd9c-54b1cfa027dc.jpeg",
      "https://a0.muscache.com/im/pictures/miso/Hosting-51573207/original/96a75a1f-408d-4c6c-9890-054c6c80fdab.jpeg",
      "https://a0.muscache.com/im/pictures/miso/Hosting-51573207/original/2c802ef1-19ef-417c-aa1e-b95deb85a885.jpeg",
      "https://a0.muscache.com/im/pictures/miso/Hosting-51573207/original/cc98001c-d734-42a2-a0c3-b0f24b888895.jpeg",
      "https://a0.muscache.com/im/pictures/hosting/Hosting-U3RheVN1cHBseUxpc3Rpbmc6NTE1NzMyMDc%3D/original/98cb244c-1f08-4b84-95ba-e714d476e35f.jpeg"
    ],
    "amenities": [
      "Air conditioning",
      "Heating",
      "Wifi",
      "Kitchen",
      "Free parking on premises",
      "Essentials",
      "Hair dryer",
      "Iron",
      "Laptop friendly workspace",
      "TV"
    ],
    "reviews": {
      "score": 4.9,
      "numberOfReviews": 10,
      "scoreCleanliness": 5,
      "scoreAccuracy": 5,
      "scoreCommunication": 5,
      "scoreLocation": 5,
      "scoreCheckIn": 5,
      "scoreValue": 5
    }
  }  
]
```

Also, check the [examples](examples) directory for integrating it as a library.


