package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	productcategoryquantityrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category_quantity"
	productcategoryquantityusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category_quantity"
)

func NewProductCategoryQuantityModule(db *bun.DB, chi server.ServerChi) (*productcategoryquantityrepositorybun.QuantityRepositoryBun, *productcategoryquantityusecases.Service, *handler.Handler) {
	repository := productcategoryquantityrepositorybun.NewQuantityRepositoryBun(db)
	service := productcategoryquantityusecases.NewService(repository)
	handler := handlerimpl.NewHandlerQuantity(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
