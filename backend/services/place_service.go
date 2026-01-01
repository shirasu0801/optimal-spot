package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"weather-map-suggest/models"
)

type PlaceService struct {
	ApiKey string
}

func NewPlaceService() *PlaceService {
	return &PlaceService{
		ApiKey: os.Getenv("GOOGLE_MAPS_API_KEY"),
	}
}

func (s *PlaceService) GetNearbySpots(lat, lng float64) ([]models.Spot, error) {
	if s.ApiKey == "" {
		// Mock Data
		return []models.Spot{
			{Name: "Mock Temple", Latitude: lat + 0.001, Longitude: lng + 0.001, Rating: 4.5, UserRatingsTotal: 1000, Types: []string{"tourist_attraction"}},
			{Name: "Mock Park", Latitude: lat - 0.001, Longitude: lng - 0.001, Rating: 4.2, UserRatingsTotal: 500, Types: []string{"park"}},
			{Name: "Mock Museum", Latitude: lat + 0.002, Longitude: lng - 0.002, Rating: 4.8, UserRatingsTotal: 200, Types: []string{"museum"}},
			{Name: "Mock Cafe", Latitude: lat - 0.002, Longitude: lng + 0.002, Rating: 3.9, UserRatingsTotal: 50, Types: []string{"cafe"}},
		}, nil
	}

	// Radius 2000 meters
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%f,%f&radius=2000&type=tourist_attraction&key=%s", lat, lng, s.ApiKey)
	
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Place API request failed: %v. Using mock data.\n", err)
		return s.getMockSpots(lat, lng), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Place API returned status: %s. Using mock data.\n", resp.Status)
		return s.getMockSpots(lat, lng), nil
	}

	var result struct {
		Results []struct {
			PlaceID          string `json:"place_id"`
			Name             string `json:"name"`
			Geometry         struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
			Rating           float64  `json:"rating"`
			UserRatingsTotal int      `json:"user_ratings_total"`
			Types            []string `json:"types"`
			Photos           []struct {
				PhotoReference string `json:"photo_reference"`
			} `json:"photos"`
		} `json:"results"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Failed to decode place response: %v. Using mock data.\n", err)
		return s.getMockSpots(lat, lng), nil
	}

	if result.Status != "OK" && result.Status != "ZERO_RESULTS" {
		fmt.Printf("Place API error status: %s. Using mock data.\n", result.Status)
		return s.getMockSpots(lat, lng), nil
	}

	var spots []models.Spot
	for _, r := range result.Results {
		photoRef := ""
		if len(r.Photos) > 0 {
			photoRef = r.Photos[0].PhotoReference
		}
		spots = append(spots, models.Spot{
			ID:               r.PlaceID,
			Name:             r.Name,
			Latitude:         r.Geometry.Location.Lat,
			Longitude:        r.Geometry.Location.Lng,
			Rating:           r.Rating,
			UserRatingsTotal: r.UserRatingsTotal,
			Types:            r.Types,
			PhotoReference:   photoRef,
		})
	}

	return spots, nil
}

func (s *PlaceService) getMockSpots(lat, lng float64) []models.Spot {
	return []models.Spot{
		{Name: "Mock Temple (API Error Fallback)", Latitude: lat + 0.001, Longitude: lng + 0.001, Rating: 4.5, UserRatingsTotal: 1000, Types: []string{"tourist_attraction"}},
		{Name: "Mock Park", Latitude: lat - 0.001, Longitude: lng - 0.001, Rating: 4.2, UserRatingsTotal: 500, Types: []string{"park"}},
		{Name: "Mock Museum", Latitude: lat + 0.002, Longitude: lng - 0.002, Rating: 4.8, UserRatingsTotal: 200, Types: []string{"museum"}},
	}
}
