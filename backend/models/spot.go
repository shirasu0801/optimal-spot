package models

type Spot struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Rating      float64 `json:"rating"`      // Google rating
	UserRatingsTotal int `json:"user_ratings_total"`
	Types       []string `json:"types"`
	PhotoReference string `json:"photo_reference,omitempty"`
	// Additional computed fields
	CrowdLevel  string  `json:"crowd_level"` // "High", "Medium", "Low"
	WeatherSuitability string `json:"weather_suitability"` // "Good", "Bad"
	Score       float64 `json:"score"` // Internal score for ranking
}

type WeatherInfo struct {
	Main        string  `json:"main"` // Clear, Rain, etc.
	Description string  `json:"description"`
	Temp        float64 `json:"temp"`
}
