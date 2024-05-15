package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	addressrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/address"
)

func NewAddressModule(db *bun.DB, chi *server.ServerChi) *addressrepositorybun.AddressRepositoryBun {
	repository := addressrepositorybun.NewAddressRepositoryBun(db)
	return repository
}
