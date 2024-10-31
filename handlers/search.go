package handlers

import (
    "encoding/json"
    "net/http"
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
                    Text:    location,
                    Context: artistName,
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
                // Check for duplicates before adding
                isDuplicate := false
                for _, existing := range allResults {
                    if existing.ID == location.ArtistId && existing.Text == loc {
                        isDuplicate = true
                        break
                    }
                }
                if !isDuplicate {
                    allResults = append(allResults, SearchResult{
                        Type:    "location",
                        ID:      location.ArtistId,
                        Text:    loc,
                        Context: artistName,
                    })
                }
            }
        }
    }

    // Limit results to 10
    if len(allResults) > 10 {
        allResults = allResults[:10]
    }

    return allResults, nil
}

// searchDates searches for dates in artist's locations
func searchDates(query string) ([]SearchResult, error) {
    var results []SearchResult
    
    for _, artist := range AllArtistInfo {
        for location, dates := range artist.DateAndLocation {
            for _, date := range dates {
                if strings.Contains(strings.ToLower(date), strings.ToLower(query)) {
                    results = append(results, SearchResult{
                        Type:    "date",
                        ID:      artist.Id,
                        Text:    date,
                        Context: location,
                    })
                }
            }
        }
        
        if len(results) >= 10 {
            break
        }
    }
    
    return results, nil
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
    
    // Define search functions while maintaining original search order
    searchFuncs := []searchFunction{
        searchArtists,
        searchLocations,
        searchDates,
    }
    
    // Perform all searches while respecting the result limit
    for _, searchFunc := range searchFuncs {
        var results []SearchResult
        results, err = searchFunc(query)
        if err != nil {
            http.Error(w, "Error performing search: "+err.Error(), http.StatusInternalServerError)
            return
        }
        
        // Add new results, avoiding duplicates
        for _, result := range results {
            isDuplicate := false
            for _, existing := range allResults {
                if existing.Type == result.Type && existing.ID == result.ID && existing.Text == result.Text {
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