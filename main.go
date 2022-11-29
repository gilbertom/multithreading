package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ApiCep struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func main() {
	chApiCep := make(chan ApiCep)
	chViaCep := make(chan ViaCep)

	go func() {
		cep, error := getApiCep("04136-031")
		if error != nil {
			panic(error)
		}
		chApiCep <- *cep
	}()

	go func() {
		cep, error := getViaCep("04136-031")
		if error != nil {
			panic(error)
		}
		chViaCep <- *cep
	}()

	select {
	case response := <-chApiCep:
		fmt.Printf("\n\tA API que retornou mais rápido foi https://cdn.apicep.com/file/apicep\n")
		fmt.Printf("\tResponse: %v\n\n", response)

	case response := <-chViaCep:
		fmt.Printf("\n\tA API que retornou mais rápido foi http://viacep.com.br/ws\n")
		fmt.Printf("\tResponse: %v\n\n", response)

	case <-time.After(time.Second * 1):
		println("\n\tTimeout Error. Exceeded one second.")
	}
}

func getViaCep(cep string) (*ViaCep, error) {
	req, error := http.Get("http://viacep.com.br/ws/" + cep + "/json")
	if error != nil {
		return nil, error
	}
	defer req.Body.Close()

	body, error := io.ReadAll(req.Body)
	if error != nil {
		return nil, error
	}

	var c ViaCep
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	return &c, nil
}

func getApiCep(cep string) (*ApiCep, error) {
	req, error := http.Get("https://cdn.apicep.com/file/apicep/" + cep + ".json")
	if error != nil {
		return nil, error
	}
	defer req.Body.Close()

	body, error := io.ReadAll(req.Body)
	if error != nil {
		return nil, error
	}

	var c ApiCep
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	return &c, nil
}
