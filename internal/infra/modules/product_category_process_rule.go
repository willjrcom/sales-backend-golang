package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	productcategoryprocessrulerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category_process_rule"
	productcategoryprocessruleusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category_process_rule"
)

func NewProductCategoryProcessRuleModule(db *bun.DB, chi server.ServerChi) (*productcategoryprocessrulerepositorybun.ProcessRuleRepositoryBun, *productcategoryprocessruleusecases.Service, *handler.Handler) {
	repository := productcategoryprocessrulerepositorybun.NewProcessRuleRepositoryBun(db)
	service := productcategoryprocessruleusecases.NewService(repository)
	handler := handlerimpl.NewHandlerProcessRuleCategory(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
