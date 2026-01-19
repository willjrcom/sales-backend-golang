package model

import (
	"github.com/uptrace/bun"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Client struct {
	entitymodel.Entity
	Person
	bun.BaseModel `bun:"table:clients"`
	ClienteCommonAttributes
}

type ClienteCommonAttributes struct {
	IsActive bool `bun:"column:is_active,type:boolean"`
}

func (c *Client) FromDomain(client *cliententity.Client) {
	if client == nil {
		return
	}
	*c = Client{
		Entity: entitymodel.FromDomain(client.Entity),
		ClienteCommonAttributes: ClienteCommonAttributes{
			IsActive: client.IsActive,
		},
	}
	c.Person.FromDomain(&client.Person)
}

func (c *Client) ToDomain() *cliententity.Client {
	if c == nil {
		return nil
	}
	return &cliententity.Client{
		Entity: c.Entity.ToDomain(),
		Person: *c.Person.ToDomain(),
		ClienteCommonAttributes: cliententity.ClienteCommonAttributes{
			IsActive: c.IsActive,
		},
	}
}
