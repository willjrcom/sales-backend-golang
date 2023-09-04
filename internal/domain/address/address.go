package addressentity

import "github.com/google/uuid"

type Address struct {
	ID         uuid.UUID
	street     string
	Number     string
	Complement string
	Reference  string
	City       string
	State      string
	Cep        string
}

func NewAddress() *Address {
	return &Address{}
}
