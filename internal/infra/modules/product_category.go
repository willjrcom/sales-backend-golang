package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	productcategoryrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category"
	productcategoryusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category"
)

func NewProductCategoryModule(db *bun.DB, chi *server.ServerChi) (model.CategoryRepository, *productcategoryusecases.Service, *handler.Handler) {
	repository := productcategoryrepositorybun.NewProductCategoryRepositoryBun(db)
	service, err := productcategoryusecases.InitializeService()
	if err != nil {
		panic(err)
	}

	handler := handlerimpl.NewHandlerProductCategory(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
