package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	sponsorrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres"
	sponsorusecases "github.com/willjrcom/sales-backend-go/internal/usecases/sponsor"
)

func NewSponsorModule(db *bun.DB, chi *server.ServerChi) (model.SponsorRepository, *sponsorusecases.SponsorService, *handler.Handler) {
	repo := sponsorrepositorybun.NewSponsorRepository(db)
	service := sponsorusecases.NewSponsorService(repo)
	handlerObj := handlerimpl.NewHandlerSponsor(service)
	chi.AddHandler(handlerObj)
	return repo, service, handlerObj
}
