package itemdto

import "errors"

var ErrNameRequired = errors.New("name is required")

type RemovedItemDTO struct {
	Name *string `json:"name"`
}

func (a *RemovedItemDTO) validate() error {
	if a.Name == nil || *a.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (a *RemovedItemDTO) ToDomain() (*string, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}

	return a.Name, nil
}
