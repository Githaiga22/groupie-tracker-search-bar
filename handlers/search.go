package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	model "tracker/models"
)

// SearchResult with additional context field
type SearchResult struct {
	Type    string `json:"type"`
	ID      int    `json:"id"`
	Text    string `json:"text"`
	Context string `json:"context,omitempty"` // Optional context like artist name
}

type SearchResponse struct {
	Success bool           `json:"success"`
	Results []SearchResult `json:"results"`
}

// searchByType handles searching for a specific type of content
type searchFunction func(query string) ([]SearchResult, error)

// getArtistNameById returns artist name for a given ID
func getArtistNameById(id int) string {
	for _, artist := range AllArtistInfo {
		if artist.Id == id {
			return artist.Name
		}
	}
	return ""
}

// getArtistNameById returns artist name for a given ID
func getArtistCreationbyId(id int) string {
	for _, artist := range AllArtistInfo {
		if artist.Id == id {
			return strconv.Itoa(artist.CreationDate)
		}
	}
	return ""
}

// searchArtists searches for artists by name
func searchArtists(query string) ([]SearchResult, error) {
	var results []SearchResult

	for _, artist := range AllArtistInfo {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			results = append(results, SearchResult{
				Type: "artist",
				ID:   artist.Id,
				Text: artist.Name,
			})
		}

		if len(results) >= 10 {
			break
		}
	}

	return results, nil
}

// searchCreations searches for creation dates
func searchCreations(query string) ([]SearchResult, error) {
	var results []SearchResult

	for _, artist := range AllArtistInfo {
		if strings.Contains(strconv.Itoa(artist.CreationDate), strings.ToLower(query)) {
			results = append(results, SearchResult{
				Type:    "creation",
				ID:      artist.Id,
				Text:    artist.Name,
				Context: strconv.Itoa(artist.CreationDate),
			})
		}

		if len(results) >= 10 {
			break
		}
	}

	return results, nil
}

// searchFirstAlbum searches for artists by First Album
func searchFirstAlbum(query string) ([]SearchResult, error) {
	var results []SearchResult

	for _, artist := range AllArtistInfo {
		if strings.Contains(artist.FirstAlbum, strings.ToLower(query)) {
			results = append(results, SearchResult{
				Type:    "First Album",
				ID:      artist.Id,
				Text:    artist.Name,
				Context: artist.FirstAlbum,
			})
		}

		if len(results) >= 10 {
			break
		}
	}

	return results, nil
}

// searchMembers searches for artists by members
func searchMembers(query string) ([]SearchResult, error) {
	var results []SearchResult

	for _, artist := range AllArtistInfo {

		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), strings.ToLower(query)) {
				results = append(results, SearchResult{
					Type:    "member",
					ID:      artist.Id,
					Text:    artist.Name,
					Context: member,
				})
			}

			if len(results) >= 10 {
				break
			}
		}

	}

	return results, nil
}

// searchLocations searches in both locations and relations endpoints
func searchLocations(query string) ([]SearchResult, error) {
	var allResults []SearchResult

	// Fetch relations data
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var relationsData model.RootsRelation
	if err := json.NewDecoder(resp.Body).Decode(&relationsData); err != nil {
		return nil, err
	}

	// Search through relations data
	for _, relation := range relationsData.Relation {
		artistName := getArtistNameById(relation.Id)
		for location := range relation.Places {
			if strings.Contains(strings.ToLower(location), strings.ToLower(query)) {
				allResults = append(allResults, SearchResult{
					Type:    "location",
					ID:      relation.Id,
					Text:    artistName,
					Context: location,
				})
			}
		}
	}

	// Now fetch and search locations data
	locResp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		return nil, err
	}
	defer locResp.Body.Close()

	var locationsData model.AllLocations
	if err := json.NewDecoder(locResp.Body).Decode(&locationsData); err != nil {
		return nil, err
	}

	// Search through locations data
	for _, location := range locationsData.Location {
		artistName := getArtistNameById(location.ArtistId)
		for _, loc := range location.Locations {
			if strings.Contains(strings.ToLower(loc), strings.ToLower(query)) {
				// Modified duplicate check to consider the actual location (Context)
				isDuplicate := false
				for _, existing := range allResults {
					if existing.ID == location.ArtistId && existing.Context == loc {
						isDuplicate = true
						break
					}
				}
				if !isDuplicate {
					allResults = append(allResults, SearchResult{
						Type:    "location",
						ID:      location.ArtistId,
						Text:    artistName,
						Context: loc,
					})
				}
			}
		}
	}

	return allResults, nil
}

// SearchHandler handles the search endpoint
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		wrongMethodHandler(w)
		return
	}

	query := strings.ToLower(r.URL.Query().Get("q"))
	if query == "" {
		json.NewEncoder(w).Encode(SearchResponse{
			Success: true,
			Results: []SearchResult{},
		})
		return
	}

	var allResults []SearchResult
	var err error

	searchFuncs := []searchFunction{
		searchArtists,
		searchLocations,
		searchCreations,
		searchFirstAlbum,
		searchMembers,
	}

	// Perform all searches while respecting the result limit
	for _, searchFunc := range searchFuncs {
		var results []SearchResult
		results, err = searchFunc(query)
		if err != nil {
			http.Error(w, "Error performing search: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Modified duplicate check to consider Type and Context
		for _, result := range results {
			isDuplicate := false
			for _, existing := range allResults {
				if existing.Type == result.Type &&
					existing.ID == result.ID &&
					existing.Context == result.Context {
					isDuplicate = true
					break
				}
			}

			if !isDuplicate {
				allResults = append(allResults, result)
			}

			if len(allResults) >= 10 {
				break
			}
		}

		if len(allResults) >= 10 {
			break
		}
	}

	// Limit results to 10
	if len(allResults) > 10 {
		allResults = allResults[:10]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SearchResponse{
		Success: true,
		Results: allResults,
	})
}
