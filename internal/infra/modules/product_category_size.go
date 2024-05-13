package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	productcategorysizerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category_size"
	productcategorysizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category_size"
)

func NewProductCategorySizeModule(db *bun.DB, chi server.ServerChi) (*productcategorysizerepositorybun.SizeRepositoryBun, *productcategorysizeusecases.Service, *handler.Handler) {
	repository := productcategorysizerepositorybun.NewSizeRepositoryBun(db)
	service := productcategorysizeusecases.NewService(repository)
	handler := handlerimpl.NewHandlerSize(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
