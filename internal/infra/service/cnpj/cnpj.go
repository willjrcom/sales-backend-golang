package cnpj

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
)

var (
	ErrCNPJNotFound = errors.New("cnpj not found")
)

const url = "https://www.receitaws.com.br/v1/cnpj/"

type Cnpj struct {
	Cnpj         string `json:"cnpj"`
	BusinessName string `json:"nome"`
	TradeName    string `json:"fantasia"`
	Street       string `json:"logradouro"`
	Number       string `json:"numero"`
	City         string `json:"municipio"`
	Neighborhood string `json:"bairro"`
	UF           string `json:"uf"`
	Cep          string `json:"cep"`
}

func Get(cnpjString string) (*Cnpj, error) {
	regexNumeros := regexp.MustCompile("[^0-9]")

	// Remover caracteres não numéricos
	cnpjNumeros := regexNumeros.ReplaceAllString(cnpjString, "")

	response, err := http.Get(url + cnpjNumeros)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Ler o corpo da response
	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var cnpj Cnpj
	if err := json.Unmarshal(body, &cnpj); err != nil {
		return nil, err
	}

	if cnpj.Cnpj == "" {
		return nil, ErrCNPJNotFound
	}

	return &cnpj, nil
}
