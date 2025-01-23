package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	productrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product"
	productusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product"
)

func NewProductModule(db *bun.DB, chi *server.ServerChi) (model.ProductRepository, *productusecases.Service, *handler.Handler) {
	repository := productrepositorybun.NewProductRepositoryBun(db)
	service, err := productusecases.InitializeService()
	if err != nil {
		panic(err)
	}

	handler := handlerimpl.NewHandlerProduct(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
