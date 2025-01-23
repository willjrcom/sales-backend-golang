package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	addressrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/address"
)

func NewAddressModule(db *bun.DB, chi *server.ServerChi) model.AddressRepository {
	repository := addressrepositorybun.NewAddressRepositoryBun(db)
	return repository
}
