package services

import (
	"strings"
	"weather-map-suggest/models"
)

type RankingService struct{}

func NewRankingService() *RankingService {
	return &RankingService{}
}

// SuggestTopSpots filters and ranks spots based on weather and estimated crowd.
// weather: Current weather info
// spots: List of nearby spots
// Returns: Top 3 sorted spots
func (s *RankingService) SuggestTopSpots(weather *models.WeatherInfo, spots []models.Spot) []models.Spot {
	for i := range spots {
		score := 0.0
		
		// Base score from rating (0-5) * 10 -> 0-50
		score += spots[i].Rating * 10

		// Weather adjustment
		isRaining := strings.Contains(strings.ToLower(weather.Main), "rain") || strings.Contains(strings.ToLower(weather.Description), "rain")
		isOutdoor := isOutdoorSpot(spots[i].Types)

		if isRaining {
			if !isOutdoor {
				score += 20 // Boost indoor spots
				spots[i].WeatherSuitability = "Good (Indoor)"
			} else {
				score -= 30 // Penalize outdoor spots heavily
				spots[i].WeatherSuitability = "Bad (Outdoor)"
			}
		} else {
			// Sunny/Cloudy
			if isOutdoor {
				score += 15 // Boost outdoor spots slightly
				spots[i].WeatherSuitability = "Good (Outdoor)"
			} else {
				spots[i].WeatherSuitability = "Neutral"
			}
		}

		// Crowd adjustment (heuristic based on UserRatingsTotal)
		// Assume > 2000 ratings means highly crowded
		if spots[i].UserRatingsTotal > 2000 {
			score -= 10 // Penalize very crowded spots
			spots[i].CrowdLevel = "High"
		} else if spots[i].UserRatingsTotal > 500 {
			spots[i].CrowdLevel = "Medium"
		} else {
			score += 5 // Boost hidden gems
			spots[i].CrowdLevel = "Low"
		}
		
		spots[i].Score = score
	}

	// Sort by Score descending
	// Simple bubble sort for small list
	n := len(spots)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if spots[j].Score < spots[j+1].Score {
				spots[j], spots[j+1] = spots[j+1], spots[j]
			}
		}
	}

	// Return top 3
	if len(spots) > 3 {
		return spots[:3]
	}
	return spots
}

func isOutdoorSpot(types []string) bool {
	outdoorTypes := []string{"park", "amusement_park", "campground", "zoo", "stadium", "tourist_attraction", "natural_feature"}
	for _, t := range types {
		for _, ot := range outdoorTypes {
			if t == ot {
				return true
			}
		}
	}
	return false
}
