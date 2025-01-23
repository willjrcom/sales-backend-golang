//go:build wireinject
// +build wireinject

package quantityusecases

import (
	"github.com/google/wire"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	categoryrepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/category"
	quantityrepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/quantity"
	categoryrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category"
	quantityrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/quantity"
)

func InitializeService() (*Service, error) {
	wire.Build(
		database.NewPostgreSQLConnection,                      // Provider para *bun.DB
		quantityrepositorybun.NewQuantityRepositoryBun,        // Provider para QuantityRepository
		categoryrepositorybun.NewProductCategoryRepositoryBun, // Provider para CategoryRepository
		NewService, // Construtor do serviço principal
	)
	return nil, nil
}

func InitializeServiceForTest() (*Service, error) {
	wire.Build(
		quantityrepositorylocal.NewQuantityRepositoryLocal, // Provider para QuantityRepository
		categoryrepositorylocal.NewCategoryRepositoryLocal, // Provider para CategoryRepository
		NewService, // Construtor do serviço principal
	)
	return nil, nil
}
