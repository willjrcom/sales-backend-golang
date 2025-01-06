package cliententity

import (
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type Client struct {
	personentity.Person
}

func NewClient(person *personentity.Person) *Client {
	return &Client{
		Person: *person,
	}
}
