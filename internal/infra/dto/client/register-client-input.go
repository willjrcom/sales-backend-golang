package clientdto

import (
	"errors"
	"time"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	persondto "github.com/willjrcom/sales-backend-go/internal/infra/dto/person"
)

type RegisterClientInput struct {
	persondto.Person `json:"person"`
}

func (r *RegisterClientInput) Validate() error {
	if r.Person.Name == "" {
		return errors.New("name is required")
	}

	if r.Person.Address.Street == "" || r.Person.Address.Number == "" || r.Person.Address.Neighborhood == "" {
		return errors.New("Address is required")
	}

	if len(r.Person.Contacts) == 0 {
		return errors.New("Contacts is required")
	}

	return nil
}

func (r *RegisterClientInput) ToModel() (*cliententity.Client, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	adress := addressentity.Address{
		Street:       r.Person.Address.Street,
		Number:       r.Person.Address.Number,
		Complement:   r.Person.Address.Complement,
		Reference:    r.Person.Address.Reference,
		Neighborhood: r.Person.Address.Neighborhood,
		City:         r.Person.Address.City,
		State:        r.Person.Address.State,
		Cep:          r.Person.Address.Cep,
	}

	person := personentity.Person{
		Name:     r.Person.Name,
		Birthday: r.Person.Birthday,
		Email:    r.Person.Email,
		Contacts: r.Person.Contacts,
		Cpf:      r.Person.Cpf,
		Address:  adress,
	}

	return &cliententity.Client{
		Person:       person,
		TotalOrders:  0,
		DateRegister: time.Now(),
	}, nil
}
