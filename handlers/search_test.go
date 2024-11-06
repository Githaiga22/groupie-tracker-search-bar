package handlers

import (
	"testing"
	"tracker/models"
)

func Test_getArtistNameById(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"valid id", args{1}, ""},
		{"valid id", args{0}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getArtistNameById(tt.args.id); got != tt.want {
				t.Errorf("getArtistNameById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getArtistCreationbyId(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"valid id", args{1}, ""},
		{"valid id", args{0}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getArtistCreationbyId(tt.args.id); got != tt.want {
				t.Errorf("getArtistCreationbyId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchArtists(t *testing.T) {
	// Set up test data
	AllArtistInfo = []models.Data{
		{
			Id:           1,
			Name:         "The Beatles",
			Members:      []string{"John", "Paul", "George", "Ringo"},
			CreationDate: 1960,
			FirstAlbum:   "Please Please Me",
		},
		{
			Id:           2,
			Name:         "Queen",
			Members:      []string{"Freddie", "Brian", "Roger", "John"},
			CreationDate: 1970,
			FirstAlbum:   "Queen",
		},
		{
			Id:           3,
			Name:         "Pink Floyd",
			Members:      []string{"Roger", "David", "Nick", "Richard"},
			CreationDate: 1965,
			FirstAlbum:   "The Piper at the Gates of Dawn",
		},
	}

	tests := []struct {
		name    string
		query   string
		wantLen int
	}{
		{"existing artist name", "beatles", 1},
		{"non-existing artist", "unknown", 0},
		{"partial match", "een", 1}, // Should match "Queen"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchArtists(tt.query)
			if err != nil {
				t.Fatalf("searchArtists() error = %v", err)
			}
			if len(got) != tt.wantLen {
				t.Errorf("searchArtists() got = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}
