package providers

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fogleman/gg"
	"github.com/nerijusdu/esp-tv-api/src/constants"
	"github.com/nerijusdu/esp-tv-api/src/util"
	"github.com/patrickmn/go-cache"
)

type WeatherProvider struct {
	ApiKey   string
	Location string
	Units    string
	cache    *cache.Cache
}

type WeatherConfig struct {
	Location string `json:"location"`
	Units    string `json:"units"`
}

func (w *WeatherProvider) Init(config any) error {
	cfg, err := util.CastConfig[WeatherConfig](config)
	if err != nil {
		return err
	}

	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("OPENWEATHERMAP_API_KEY environment variable is not set")
	}
	w.ApiKey = apiKey
	w.Location = cfg.Location

	w.Units = "metric"
	if cfg.Units != "" {
		w.Units = cfg.Units
	}

	w.cache = cache.New(30*time.Minute, 1*time.Hour)

	return nil
}

func (w *WeatherProvider) GetView(cursor string) (ViewResponse, error) {
	weatherData, err := w.getWeatherData()
	if err != nil {
		return ViewResponse{}, err
	}

	img, err := w.renderWeather(weatherData)
	if err != nil {
		return ViewResponse{}, err
	}

	return ViewResponse{
		View: View{
			Data:         img,
			RefreshAfter: 10000,
		},
		Cursor:     cursor,
		NextCursor: "",
	}, nil
}

func (w *WeatherProvider) getWeatherData() (WeatherResponse, error) {
	cacheKey := fmt.Sprintf("weather-%s-%s", w.Location, w.Units)
	if cachedData, found := w.cache.Get(cacheKey); found {
		return cachedData.(WeatherResponse), nil
	}

	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=%s",
		w.Location,
		w.ApiKey,
		w.Units,
	)

	resp, err := http.Get(url)
	if err != nil {
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WeatherResponse{}, fmt.Errorf("weather API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherResponse{}, err
	}

	var weatherResp WeatherResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		return WeatherResponse{}, err
	}

	w.cache.Set(cacheKey, weatherResp, cache.DefaultExpiration)

	return weatherResp, nil
}

func (w *WeatherProvider) renderWeather(weather WeatherResponse) ([]byte, error) {
	dc := gg.NewContext(constants.DISPLAY_WIDTH, constants.DISPLAY_HEIGHT)
	dc.SetColor(color.White)

	dc.DrawString(weather.Name, 5, 10)

	currentTime := time.Now().Format("15:04")
	dc.DrawStringAnchored(currentTime, float64(constants.DISPLAY_WIDTH-5), 10, 1.0, 0)

	tempUnit := "°C"
	if w.Units == "imperial" {
		tempUnit = "°F"
	}

	temperature := fmt.Sprintf("%.1f%s", weather.Main.Temp, tempUnit)
	dc.DrawStringAnchored(temperature, float64(constants.DISPLAY_WIDTH/2), 20, 0.5, 0.5)

	var description string
	if len(weather.Weather) > 0 {
		description = weather.Weather[0].Description
	} else {
		description = "Unknown"
	}

	var weatherInfo string
	if w.Units == "imperial" {
		weatherInfo = fmt.Sprintf("H:%d%% W:%.1fmph", weather.Main.Humidity, weather.Wind.Speed)
	} else {
		weatherInfo = fmt.Sprintf("H:%d%% W:%.1fm/s", weather.Main.Humidity, weather.Wind.Speed)
	}

	dc.DrawStringAnchored(description, float64(constants.DISPLAY_WIDTH/2), 35, 0.5, 0.5)
	dc.DrawStringAnchored(weatherInfo, float64(constants.DISPLAY_WIDTH/2), 50, 0.5, 0.5)
	dc.Stroke()

	res := util.GraphicToBytes(dc)
	return *res, nil
}
