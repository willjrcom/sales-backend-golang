package itemdto

import "errors"

var ErrNameRequired = errors.New("name is required")

type RemovedItemInput struct {
	Name *string `json:"name"`
}

func (a *RemovedItemInput) validate() error {
	if a.Name == nil || *a.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (a *RemovedItemInput) ToModel() (*string, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}

	return a.Name, nil
}
