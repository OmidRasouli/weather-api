package openweather

import (
	"context"
	"fmt"
	"time"

	"github.com/OmidRasouli/weather-api/internal/application/interfaces"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	apiKey string
	client *resty.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: resty.New().
			SetTimeout(5 * time.Second),
	}
}

// Sample API response struct
type apiResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Dt int64 `json:"dt"` // Unix timestamp
}

func (c *Client) FetchWeatherData(ctx context.Context, city string, country string) (*interfaces.WeatherAPIResponse, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s,%s&appid=%s&units=metric", city, country, c.apiKey)

	var res apiResponse
	_, err := c.client.R().
		SetContext(ctx).
		SetResult(&res).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to call weather API: %w", err)
	}

	if len(res.Weather) == 0 {
		return nil, fmt.Errorf("invalid response: missing weather description")
	}

	return &interfaces.WeatherAPIResponse{
		Temperature: res.Main.Temp,
		Description: res.Weather[0].Description,
		Humidity:    res.Main.Humidity,
		WindSpeed:   res.Wind.Speed,
		FetchedAt:   time.Unix(res.Dt, 0),
	}, nil
}
