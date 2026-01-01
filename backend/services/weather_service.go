package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"weather-map-suggest/models"
)

type WeatherService struct {
	ApiKey string
}

func NewWeatherService() *WeatherService {
	return &WeatherService{
		ApiKey: os.Getenv("OPENWEATHER_API_KEY"),
	}
}

func (s *WeatherService) GetCurrentWeather(lat, lng float64) (*models.WeatherInfo, error) {
	if s.ApiKey == "" {
		// Mock Data
		return &models.WeatherInfo{
			Main:        "Clear",
			Description: "clear sky",
			Temp:        25.0,
		}, nil
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", lat, lng, s.ApiKey)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Weather API request failed: %v. Using mock data.\n", err)
		return s.getMockWeather(), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Weather API returned status: %s. Using mock data.\n", resp.Status)
		return s.getMockWeather(), nil
	}

	var result struct {
		Weather []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
		} `json:"weather"`
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Failed to decode weather response: %v. Using mock data.\n", err)
		return s.getMockWeather(), nil
	}

	if len(result.Weather) == 0 {
		return s.getMockWeather(), nil
	}

	return &models.WeatherInfo{
		Main:        result.Weather[0].Main,
		Description: result.Weather[0].Description,
		Temp:        result.Main.Temp,
	}, nil
}

func (s *WeatherService) getMockWeather() *models.WeatherInfo {
	return &models.WeatherInfo{
		Main:        "Fine (Mock)",
		Description: "sunny (mock data due to api error)",
		Temp:        22.0,
	}
}
