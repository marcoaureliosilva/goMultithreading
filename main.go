package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// Estrutura para armazenar a resposta da API
type Address struct {
	Cep          string `json:"cep"`
	Logradouro   string `json:"logradouro"`
	Street       string `json:"street"`
	Complemento  string `json:"complemento"`
	Bairro       string `json:"bairro"`
	Neighborhood string `json:"neighborhood"`
	Localidade   string `json:"localidade"`
	City         string `json:"city"`
	Uf           string `json:"uf"`
	State        string `json:"state"`
	Unidade      string `json:"unidade"`
	Ibge         string `json:"ibge"`
	Gia          string `json:"gia"`
}

func main() {
	cep := "01153000"

	// Inicializando as URLs das APIs
	api1 := "https://brasilapi.com.br/api/cep/v1/" + cep
	api2 := "http://viacep.com.br/ws/" + cep + "/json/"

	// Canal para receber a resposta mais rápida
	result := make(chan string, 2)
	var wg sync.WaitGroup

	// Função anônima para gerenciar o timeout de 1 segundo
	go func() {
		wg.Wait()
		close(result)
	}()

	wg.Add(2)
	go requestAPI(api1, "BrasilAPI", &wg, result)
	go requestAPI(api2, "ViaCEP", &wg, result)

	select {
	case res := <-result:
		fmt.Println(res)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout: Nenhuma resposta foi recebida dentro do limite de 1 segundo")
	}
}

func requestAPI(url, apiName string, wg *sync.WaitGroup, result chan<- string) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		result <- fmt.Sprintf("%s: Erro ao fazer a requisição HTTP", apiName)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		result <- fmt.Sprintf("%s: Erro ao ler a resposta", apiName)
		return
	}

	var address Address
	if err := json.Unmarshal(body, &address); err != nil {
		result <- fmt.Sprintf("%s: Erro ao fazer o parse da resposta JSON", apiName)
		return
	}

	if apiName == "BrasilAPI" {
		result <- fmt.Sprintf("Resposta da %s\nRua:    %s\nBairro: %s\nCidade: %s - %s\nCEP:    %s\nComplemento: %s", apiName,
			address.Street, address.Neighborhood, address.City, address.State, address.Cep, address.Complemento)
	} else {
		result <- fmt.Sprintf("Resposta da %s\nRua:    %s\nBairro: %s\nCidade: %s - %s\nCEP:    %s\nComplemento: %s", apiName,
			address.Logradouro, address.Bairro, address.Localidade, address.Uf, address.Cep, address.Complemento)
	}
}
