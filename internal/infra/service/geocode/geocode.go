package geocodeservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
)

type GeocodeResponse struct {
	Results []struct {
		Geometry struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"geometry"`
	} `json:"results"`
	Status map[string]interface{} `json:"status"`
}

func GetCoordinates(address *addressentity.AddressCommonAttributes) (*addressentity.Coordinates, error) {
	baseURL := "https://api.opencagedata.com/geocode/v1/json"
	apiKey := "4f638e9ecf444f81b5eb379318936d87"

	// Codifica o endereço na URL
	params := url.Values{}
	params.Add("q", address.Street+" "+address.Number+", "+address.Neighborhood+", "+address.City+", "+address.State)
	params.Add("key", apiKey)

	// Constrói a URL completa
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Faz a requisição HTTP
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	// Decodifica a resposta JSON
	var geocodeResponse GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&geocodeResponse); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta JSON: %v", err)
	}

	// Verifica o status da resposta
	status, ok := geocodeResponse.Status["code"]
	if !ok || status != float64(200) {
		return nil, fmt.Errorf("erro na geocodificação: status %s", geocodeResponse.Status)
	}

	// Extrai latitude e longitude do primeiro resultado
	if len(geocodeResponse.Results) == 0 {
		return nil, fmt.Errorf("nenhum resultado encontrado para o endereço")
	}
	location := geocodeResponse.Results[0].Geometry

	return &addressentity.Coordinates{Latitude: location.Lat, Longitude: location.Lng}, nil
}
