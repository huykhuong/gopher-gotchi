package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type WeatherResponse struct {
	Current struct {
		WeatherCode int `json:"weather_code"`
		IsDay       int `json:"is_day"`
	} `json:"current"`
}

type ProcessedWeatherData struct {
	IsDay     bool
	Condition string
}

func GetWeatherData() (ProcessedWeatherData, error) {
	weatherResponse, err := fetchWeather()
	if err != nil {
		return ProcessedWeatherData{}, err
	}

	return ProcessedWeatherData{
		IsDay:     weatherResponse.Current.IsDay == 1,
		Condition: getWeatherCondition(weatherResponse.Current.WeatherCode),
	}, nil
}

func getWeatherCondition(code int) string {
	switch code {
	case 0, 1:
		return "Sunny / Clear"
	case 2, 3:
		return "Cloudy"
	case 51, 53, 55, 61, 63, 65, 80, 81, 82:
		return "Raining"
	case 95, 96, 99:
		return "Stormy"
	default:
		return "Unspecified / Foggy"
	}
}

func fetchWeather() (WeatherResponse, error) {
	basePath := "https://api.open-meteo.com/v1/forecast"
	params := url.Values{}
	params.Add("latitude", "10.823")
	params.Add("longitude", "106.6296")
	params.Add("current", "weather_code,is_day")
	params.Add("forecast_days", "1")

	fullURL := fmt.Sprintf("%s?%s", basePath, params.Encode())

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(fullURL)
	if err != nil {
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WeatherResponse{}, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var weatherResponse WeatherResponse

	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return WeatherResponse{}, err
	}

	return weatherResponse, nil
}
