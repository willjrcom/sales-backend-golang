package orderpickupdto

type UpdateOrderPickupInput struct {
	Name string `json:"name"`
}

func (o *UpdateOrderPickupInput) validate() error {
	if o.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (o *UpdateOrderPickupInput) ToDomain() (string, error) {
	if err := o.validate(); err != nil {
		return "", err
	}

	return o.Name, nil
}
