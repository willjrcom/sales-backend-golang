package addressentity

type Repository interface {
	RegisterAddress(address *Address) error
	UpdateAddress(address *Address) error
	RemoveAddress(id string) error
}
