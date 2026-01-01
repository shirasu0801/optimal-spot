package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"weather-map-suggest/services"
)

type SuggestRequest struct {
	Lat float64 `json:"lat" binding:"required"`
	Lng float64 `json:"lng" binding:"required"`
}

func main() {
	// Load .env file if exists
	godotenv.Load()

	r := gin.Default()

	// Initialize services
	weatherService := services.NewWeatherService()
	placeService := services.NewPlaceService()
	rankingService := services.NewRankingService()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.POST("/api/suggest", func(c *gin.Context) {
		var req SuggestRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 1. Get Weather
		weather, err := weatherService.GetCurrentWeather(req.Lat, req.Lng)
		if err != nil {
			// In case of error (or missing API key logic failure), continue with dummy weather or fail?
			// For robustness, logging error and proceeding might be better, but let's fail for now to verify integration.
			// Actually, let's log and use a default dummy weather if fail, so the app is usable.
			log.Printf("Weather API error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weather data"})
			return
		}

		// 2. Get Nearby Spots
		spots, err := placeService.GetNearbySpots(req.Lat, req.Lng)
		if err != nil {
			log.Printf("Place API error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch nearby spots"})
			return
		}

		// 3. Rank and Filter
		topSpots := rankingService.SuggestTopSpots(weather, spots)

		c.JSON(http.StatusOK, gin.H{
			"weather": weather,
			"suggestions": topSpots,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
