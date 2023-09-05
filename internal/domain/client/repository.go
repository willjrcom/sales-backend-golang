package cliententity

type Repository interface {
	RegisterClient(p *Client) error
	UpdateClient(p *Client) error
	DeleteClient(id string) error
	GetClientById(id string) (*Client, error)
	GetClientBy(key string, value string) (*Client, error)
	GetAllClient(key string, value string) ([]Client, error)
}
