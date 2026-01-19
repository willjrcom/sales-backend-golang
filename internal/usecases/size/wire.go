//go:build wireinject
// +build wireinject

package sizeusecases

import (
	"github.com/google/wire"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	categoryrepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/product_category"
	sizerepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/size"
	productcategoryrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category"
	sizerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/size"
)

func InitializeService() (*Service, error) {
	wire.Build(
		database.NewPostgreSQLConnection,                             // Provider para *bun.DB
		sizerepositorybun.NewSizeRepositoryBun,                       // Provider para SizeRepository
		productcategoryrepositorybun.NewProductCategoryRepositoryBun, // Provider para CategoryRepository
		NewService, // Construtor do serviço principal
	)
	return nil, nil
}

func InitializeServiceForTest() (*Service, error) {
	wire.Build(
		sizerepositorylocal.NewSizeRepositoryLocal,         // Provider para SizeRepository
		categoryrepositorylocal.NewCategoryRepositoryLocal, // Provider para CategoryRepository
		NewService, // Construtor do serviço principal
	)
	return nil, nil
}
