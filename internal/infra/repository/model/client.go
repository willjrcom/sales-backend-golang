package model

import (
	"github.com/uptrace/bun"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Client struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:clients"`
	Person
}

func (c *Client) FromDomain(client *cliententity.Client) {
	*c = Client{
		Entity: entitymodel.FromDomain(client.Entity),
	}
	c.Person.FromDomain(&client.Person)
}

func (c *Client) ToDomain() *cliententity.Client {
	return &cliententity.Client{
		Entity: c.Entity.ToDomain(),
		Person: *c.Person.ToDomain(),
	}
}
