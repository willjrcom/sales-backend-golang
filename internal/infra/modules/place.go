package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	placerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/place"
	placeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/place"
)

func NewPlaceModule(db *bun.DB, chi *server.ServerChi) (model.PlaceRepository, *placeusecases.Service, *handler.Handler) {
	repository := placerepositorybun.NewPlaceRepositoryBun(db)
	service := placeusecases.NewService(repository)
	handler := handlerimpl.NewHandlerPlace(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
