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
	r.Get("/cep/{cep}", cepHandler)

	http.ListenAndServe(":8080", r)
}

type CepResponse struct {
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

	// Quick terminal print options:
	fmt.Printf("%+v\n", viaCepResponse)

	w.WriteHeader(http.StatusOK)
	return

}

func fetchCep(cep string) (ViaCepResponse, error) {
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
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
