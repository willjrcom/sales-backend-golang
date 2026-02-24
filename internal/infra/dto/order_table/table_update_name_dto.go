package ordertabledto

import "errors"

var ErrNameRequired = errors.New("name is required")

type UpdateOrderTableNameDTO struct {
	Name string `json:"name"`
}

func (o *UpdateOrderTableNameDTO) validate() error {
	if o.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (o *UpdateOrderTableNameDTO) ToDomain() (string, error) {
	if err := o.validate(); err != nil {
		return "", err
	}

	return o.Name, nil
}
