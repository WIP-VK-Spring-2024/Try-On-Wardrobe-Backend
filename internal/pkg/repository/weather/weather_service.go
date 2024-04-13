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

func (w WeatherService) CurrentWeather(request domain.WeatherRequest) (*domain.Weather, error) {
	var queryParams url.Values
	queryParams.Add("key", w.apiKey)

	if math.Abs(float64(request.Lat)) < eps && math.Abs(float64(request.Lon)) < eps {
		queryParams.Add("q", request.IP)
	} else {
		queryParams.Add("q", fmt.Sprintf("%f,%f", request.Lat, request.Lon))
	}

	path := apiEndpoint + "/current.json?" + queryParams.Encode()
	fmt.Println("Sending request to:", path)

	resp, err := http.DefaultClient.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	weather := &domain.Weather{}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := easyjson.Unmarshal(bytes, weather); err != nil {
		return nil, err
	}

	return weather, nil
}
