package handlers

import (
    "encoding/json"
    "net/http"
    "strings"
    model "tracker/models"
)

type SearchResult struct {
    Type string `json:"type"`
    ID   int    `json:"id"`
    Text string `json:"text"`
}

type SearchResponse struct {
    Success bool           `json:"success"`
    Results []SearchResult `json:"results"`
}

// searchByType handles searching for a specific type of content
type searchFunction func(artist model.Data, query string, results []SearchResult) []SearchResult

// Search implementations for each type
func searchArtists(artist model.Data, query string, allResults []SearchResult) []SearchResult {
    if strings.Contains(strings.ToLower(artist.Name), query) {
        return append(allResults, SearchResult{
            Type: "artist",
            ID:   artist.Id,
            Text: artist.Name,
        })
    }
    return allResults
}

func searchLocations(artist model.Data, query string, allResults []SearchResult) []SearchResult {
    for location := range artist.DateAndLocation {
        if strings.Contains(strings.ToLower(location), query) {
            allResults = append(allResults, SearchResult{
                Type: "location",
                ID:   artist.Id,
                Text: location,
            })
        }
    }
    return allResults
}

func searchDates(artist model.Data, query string, allResults []SearchResult) []SearchResult {
    for _, dates := range artist.DateAndLocation {
        for _, date := range dates {
            if strings.Contains(strings.ToLower(date), query) {
                allResults = append(allResults, SearchResult{
                    Type: "date",
                    ID:   artist.Id,
                    Text: date,
                })
            }
        }
    }
    return allResults
}

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
    
    // Define search functions while maintaining original search order
    searchFuncs := []searchFunction{
        searchArtists,
        searchLocations,
        searchDates,
    }

    // Perform all searches while respecting the result limit
    for _, artist := range AllArtistInfo {
        for _, searchFunc := range searchFuncs {
            allResults = searchFunc(artist, query, allResults)
            if len(allResults) > 10 {
                allResults = allResults[:10]
                break
            }
        }
        if len(allResults) >= 10 {
            break
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(SearchResponse{
        Success: true,
        Results: allResults,
    })
}