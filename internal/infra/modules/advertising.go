package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	advertisingrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres"
	advertisingusecases "github.com/willjrcom/sales-backend-go/internal/usecases/advertising"
)

func NewAdvertisingModule(db *bun.DB, chi *server.ServerChi, sponsorRepo model.SponsorRepository, userRepo model.UserRepository) (model.AdvertisingRepository, *advertisingusecases.AdvertisingService, *handler.Handler) {
	repo := advertisingrepositorybun.NewAdvertisingRepository(db)
	service := advertisingusecases.NewAdvertisingService(repo, sponsorRepo, userRepo)
	handlerObj := handlerimpl.NewHandlerAdvertising(service)
	chi.AddHandler(handlerObj)
	return repo, service, handlerObj
}
