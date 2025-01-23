package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	schemarepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/schema"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
)

func NewSchemaModule(db *bun.DB, chi *server.ServerChi) (model.SchemaRepository, *headerservice.Service) {
	repository := schemarepositorybun.NewSchemaRepositoryBun(db)
	service := headerservice.NewService(repository)
	return repository, service
}
