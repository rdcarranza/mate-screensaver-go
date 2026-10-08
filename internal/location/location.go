package location

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Info struct {
	City      string
	Region    string
	Country   string
	Latitude  float64
	Longitude float64
	Timezone  string
}

func Now() (Info, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	response, err := client.Get("https://ipwho.is/")
	if err != nil {
		return Info{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Info{}, fmt.Errorf(
			"location service returned HTTP status %d",
			response.StatusCode,
		)
	}

	var data struct {
		Success   bool    `json:"success"`
		Message   string  `json:"message"`
		City      string  `json:"city"`
		Region    string  `json:"region"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timezone  struct {
			ID string `json:"id"`
		} `json:"timezone"`
	}

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return Info{}, err
	}

	if !data.Success {
		return Info{}, fmt.Errorf(
			"location service error: %s",
			data.Message,
		)
	}

	region := strings.TrimSuffix(data.Region, " Province")

	return Info{
		City:      data.City,
		Region:    region,
		Country:   data.Country,
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
		Timezone:  data.Timezone.ID,
	}, nil
}

func FromCoordinates(latitude, longitude float64) (Info, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	params := url.Values{}

	params.Set(
		"lat",
		strconv.FormatFloat(latitude, 'f', 6, 64),
	)

	params.Set(
		"lon",
		strconv.FormatFloat(longitude, 'f', 6, 64),
	)

	params.Set("format", "json")
	params.Set("addressdetails", "1")
	params.Set("accept-language", "es")

	endpoint := "https://nominatim.openstreetmap.org/reverse?" + params.Encode()

	request, err := http.NewRequest(
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return Info{}, err
	}

	request.Header.Set(
		"User-Agent",
		"mate-screensaver-go",
	)

	response, err := client.Do(request)
	if err != nil {
		return Info{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Info{}, fmt.Errorf(
			"reverse geocoding service returned HTTP status %d",
			response.StatusCode,
		)
	}

	var data struct {
		DisplayName string `json:"display_name"`

		Address struct {
			City         string `json:"city"`
			Town         string `json:"town"`
			Village      string `json:"village"`
			Municipality string `json:"municipality"`

			State   string `json:"state"`
			Country string `json:"country"`
		} `json:"address"`
	}

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return Info{}, err
	}

	city := data.Address.City

	if city == "" {
		city = data.Address.Town
	}

	if city == "" {
		city = data.Address.Village
	}

	if city == "" {
		city = data.Address.Municipality
	}

	return Info{
		City:      normalizeCity(city),
		Region:    data.Address.State,
		Country:   data.Address.Country,
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}

func normalizeCity(city string) string {
	const prefix = "municipio de "

	if len(city) >= len(prefix) &&
		strings.EqualFold(city[:len(prefix)], prefix) {
		return strings.TrimSpace(city[len(prefix):])
	}

	return strings.TrimSpace(city)
}
