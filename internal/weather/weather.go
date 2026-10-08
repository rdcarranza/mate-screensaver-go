package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Info struct {
	Temperature         float64
	ApparentTemperature float64
	RelativeHumidity    float64
	WindSpeed           float64
	WeatherCode         int
	IsDay               bool
	Description         string
}

func Now(latitude, longitude float64) (Info, error) {
	params := url.Values{}

	params.Set("latitude", strconv.FormatFloat(latitude, 'f', 6, 64))
	params.Set("longitude", strconv.FormatFloat(longitude, 'f', 6, 64))
	params.Set(
		"current",
		"temperature_2m,apparent_temperature,relative_humidity_2m,wind_speed_10m,weather_code,is_day",
	)

	endpoint := "https://api.open-meteo.com/v1/forecast?" + params.Encode()

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	response, err := client.Get(endpoint)
	if err != nil {
		return Info{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Info{}, fmt.Errorf(
			"weather service returned HTTP status %d",
			response.StatusCode,
		)
	}

	var data struct {
		Current struct {
			Temperature         float64 `json:"temperature_2m"`
			ApparentTemperature float64 `json:"apparent_temperature"`
			RelativeHumidity    float64 `json:"relative_humidity_2m"`
			WindSpeed           float64 `json:"wind_speed_10m"`
			WeatherCode         int     `json:"weather_code"`
			IsDay               int     `json:"is_day"`
		} `json:"current"`
	}

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return Info{}, err
	}

	return Info{
		Temperature:         data.Current.Temperature,
		ApparentTemperature: data.Current.ApparentTemperature,
		RelativeHumidity:    data.Current.RelativeHumidity,
		WindSpeed:           data.Current.WindSpeed,
		WeatherCode:         data.Current.WeatherCode,
		IsDay:               data.Current.IsDay == 1,
		Description:         description(data.Current.WeatherCode),
	}, nil
}

func description(code int) string {
	switch code {
	case 0:
		return "Despejado"

	case 1:
		return "Principalmente despejado"

	case 2:
		return "Parcialmente nublado"

	case 3:
		return "Nublado"

	case 45, 48:
		return "Niebla"

	case 51, 53, 55:
		return "Llovizna"

	case 56, 57:
		return "Llovizna helada"

	case 61, 63, 65:
		return "Lluvia"

	case 66, 67:
		return "Lluvia helada"

	case 71, 73, 75, 77:
		return "Nieve"

	case 80, 81, 82:
		return "Chaparrones"

	case 85, 86:
		return "Chaparrones de nieve"

	case 95:
		return "Tormenta"

	case 96, 99:
		return "Tormenta con granizo"

	default:
		return "Condiciones desconocidas"
	}
}
