package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	processrulerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/process_rule"
	processruleusecases "github.com/willjrcom/sales-backend-go/internal/usecases/process_rule"
)

func NewProductCategoryProcessRuleModule(db *bun.DB, chi *server.ServerChi) (*processrulerepositorybun.ProcessRuleRepositoryBun, *processruleusecases.Service, *handler.Handler) {
	repository := processrulerepositorybun.NewProcessRuleRepositoryBun(db)
	service := processruleusecases.NewService(repository)
	handler := handlerimpl.NewHandlerProcessRuleCategory(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
