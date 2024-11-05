package handlers

import (
	"testing"
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
