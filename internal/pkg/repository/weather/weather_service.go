package weather

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"

	"try-on/internal/pkg/domain"

	"github.com/mailru/easyjson"
)

const apiEndpoint = "https://api.weatherapi.com/v1"

const eps = 1e-7

type WeatherService struct {
	apiKey string
}

func New(apiKey string) domain.WeatherService {
	return &WeatherService{
		apiKey: apiKey,
	}
}

//easyjson:json
type weatherApiResponse struct {
	Location domain.GeoPosition
	Current  domain.Weather
}

func (w WeatherService) CurrentWeather(request domain.WeatherRequest) (*domain.Weather, error) {
	queryParams := make(url.Values, 2)

	if math.Abs(float64(request.Lat)) < eps && math.Abs(float64(request.Lon)) < eps {
		resp, err := w.getGeoPosition(request.IP)
		if err != nil {
			return nil, err
		}

		if math.Abs(float64(resp.Current.Temp-resp.Current.TempFahrenheit)) > eps {
			return &resp.Current, nil
		}

		request.GeoPosition = resp.Location
	}

	queryParams.Add("q", fmt.Sprintf("%f,%f", request.Lat, request.Lon))

	resp, err := w.makeRequest(queryParams)
	if err != nil {
		return nil, err
	}

	return &resp.Current, nil
}

func (w WeatherService) getGeoPosition(ip string) (*weatherApiResponse, error) {
	queryParams := make(url.Values, 2)
	queryParams.Add("q", ip)

	resp, err := w.makeRequest(queryParams)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (w WeatherService) makeRequest(params url.Values) (*weatherApiResponse, error) {
	params.Add("key", w.apiKey)

	path := apiEndpoint + "/current.json?" + params.Encode()
	fmt.Println("Sending request to:", path)

	resp, err := http.DefaultClient.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	weather := &weatherApiResponse{}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Got from weather service: ", string(bytes))

	if err := easyjson.Unmarshal(bytes, weather); err != nil {
		return nil, err
	}

	return weather, nil
}
