package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	productcategorysizerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/size"
	productcategorysizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/size"
)

func NewProductCategorySizeModule(db *bun.DB, chi *server.ServerChi) (model.SizeRepository, *productcategorysizeusecases.Service, *handler.Handler) {
	repository := productcategorysizerepositorybun.NewSizeRepositoryBun(db)
	service := productcategorysizeusecases.NewService(repository)
	handler := handlerimpl.NewHandlerSize(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
