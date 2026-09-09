package airbnb_test

import (
	"testing"

	"github.com/mikemyl/airbnb-downloader/pkg/airbnb"
)

func TestPhotoURLsFromPageState_ReturnsPhotoTourInOrder(t *testing.T) {
	state := `{"niobeClientData":[[null,{"data":{"sections":[
		{"section":{"__typename":"HeroSection","mediaItems":[{"baseUrl":"https://a0.muscache.com/im/pictures/hero.jpeg"}]}},
		{"section":{"__typename":"PhotoTourModalSection","mediaItems":[
			{"id":"1","baseUrl":"https://a0.muscache.com/im/pictures/first.jpeg"},
			{"id":"2","baseUrl":"https://a0.muscache.com/im/pictures/second.jpeg"},
			{"id":"3","baseUrl":""}
		]}}
	]}}]]}`

	photos, err := airbnb.PhotoURLsFromPageState([]byte(state))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		"https://a0.muscache.com/im/pictures/first.jpeg",
		"https://a0.muscache.com/im/pictures/second.jpeg",
	}
	if len(photos) != len(want) {
		t.Fatalf("got %d photos, want %d", len(photos), len(want))
	}
	for i, w := range want {
		if photos[i].String() != w {
			t.Errorf("photo %d = %s, want %s", i, photos[i], w)
		}
	}
}

func TestPhotoURLsFromPageState_NoPhotoTour(t *testing.T) {
	_, err := airbnb.PhotoURLsFromPageState([]byte(`{"sections":[{"__typename":"HeroSection"}]}`))
	if err == nil {
		t.Fatal("expected an error when the state has no photo tour section")
	}
}

func TestPhotoURLsFromPageState_InvalidJSON(t *testing.T) {
	_, err := airbnb.PhotoURLsFromPageState([]byte(`{not json`))
	if err == nil {
		t.Fatal("expected a decode error")
	}
}
