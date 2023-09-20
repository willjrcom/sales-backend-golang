package addressdto

type AddressOutput struct {
	ID           string `json:"id"`
	Street       string `json:"street"`
	Number       string `json:"number"`
	Complement   string `json:"complement"`
	Reference    string `json:"reference"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Cep          string `json:"cep"`
}
