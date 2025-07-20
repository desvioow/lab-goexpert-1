package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"regexp"
)

func main() {
	r := chi.NewRouter()
	r.Get("/{cep}", cepHandler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	http.ListenAndServe(":8080", r)
}

type TemperatureResponse struct {
	TempC string `json:"temp_C"`
	TempF string `json:"temp_F"`
	TempK string `json:"temp_K"`
}

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type WeatherApiResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		WindchillC float64 `json:"windchill_c"`
		WindchillF float64 `json:"windchill_f"`
		HeatindexC float64 `json:"heatindex_c"`
		HeatindexF float64 `json:"heatindex_f"`
		DewpointC  float64 `json:"dewpoint_c"`
		DewpointF  float64 `json:"dewpoint_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	regex := regexp.MustCompile(`^\d{8}$`)

	if regex.MatchString(cep) == false {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	viaCepResponse, err := fetchCep(cep)
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isViaCepResponseEmpty(viaCepResponse) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find zipcode"))
		return
	}

	weatherApiResponse, err := fetchWeather(viaCepResponse.Localidade)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find weather for the city"))
		return
	}
	if isWeatherApiResponseEmpty(weatherApiResponse) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find city"))
		return
	}

	temperature, err := getTemperatureFromWeatherApiResponse(weatherApiResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(temperature)
	w.WriteHeader(http.StatusOK)
	return

}

func fetchCep(cep string) (ViaCepResponse, error) {

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	resp, err := http.Get(url)
	if err != nil {
		return ViaCepResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ViaCepResponse{}, err
	}
	var viaCepResponse ViaCepResponse
	err = json.Unmarshal(body, &viaCepResponse)
	if err != nil {
		return ViaCepResponse{}, err
	}

	return viaCepResponse, nil
}

func isViaCepResponseEmpty(response ViaCepResponse) bool {
	empty := ViaCepResponse{}
	return response == empty
}

func fetchWeather(cityName string) (WeatherApiResponse, error) {
	const weatherApiKey = "4d0046deef4e4342bd9192050251307"
	const baseUrl = "http://api.weatherapi.com/v1"

	url := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no",
		baseUrl,
		weatherApiKey,
		cityName,
	)

	resp, err := http.Get(url)
	if err != nil {
		return WeatherApiResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherApiResponse{}, err
	}

	var weatherApiResponse WeatherApiResponse
	err = json.Unmarshal(body, &weatherApiResponse)
	if err != nil {
		return WeatherApiResponse{}, err
	}

	return weatherApiResponse, nil
}

func isWeatherApiResponseEmpty(response WeatherApiResponse) bool {
	empty := WeatherApiResponse{}
	return response == empty
}

func getTemperatureFromWeatherApiResponse(weatherApiResponse WeatherApiResponse) (TemperatureResponse, error) {

	tempC := weatherApiResponse.Current.TempC
	tempF := (tempC * 1.8) + 32
	tempK := tempC + 273

	return TemperatureResponse{
		TempC: fmt.Sprintf("%.1f", tempC),
		TempF: fmt.Sprintf("%.1f", tempF),
		TempK: fmt.Sprintf("%.1f", tempK),
	}, nil
}
