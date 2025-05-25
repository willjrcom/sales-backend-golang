package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"

	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/report"
	reportusecases "github.com/willjrcom/sales-backend-go/internal/usecases/report"
)

// NewReportModule registers report endpoints and services.
func NewReportModule(db *bun.DB, chi *server.ServerChi) {
	// Initialize core report service
	reportSvc := report.NewReportRepository(db)
	// Wrap in usecase
	usecase := reportusecases.NewService(reportSvc)
	// Create HTTP handler
	handler := handlerimpl.NewHandlerReport(usecase)
	// Register handler
	chi.AddHandler(handler)
}
