package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	productcategoryproductrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category_product"
	productcategoryproductusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category_product"
)

func NewProductCategoryProductModule(db *bun.DB, chi server.ServerChi) (*productcategoryproductrepositorybun.ProductRepositoryBun, *productcategoryproductusecases.Service, *handler.Handler) {
	repository := productcategoryproductrepositorybun.NewProductRepositoryBun(db)
	service := productcategoryproductusecases.NewService(repository)
	handler := handlerimpl.NewHandlerProduct(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
