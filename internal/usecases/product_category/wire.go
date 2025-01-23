//go:build wireinject
// +build wireinject

package productcategoryusecases

import (
	"github.com/google/wire"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	categoryrepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/category"
	productcategoryrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category"
	quantityusecases "github.com/willjrcom/sales-backend-go/internal/usecases/quantity"
	sizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/size"
)

func InitializeService() (*Service, error) {
	wire.Build(
		database.NewPostgreSQLConnection,                             // Provider para *bun.DB
		productcategoryrepositorybun.NewProductCategoryRepositoryBun, // Provider para CategoryRepository
		quantityusecases.InitializeService,                           // Provider para QuantityRepository
		sizeusecases.InitializeService,
		NewService, // Construtor do serviço principal
	)
	return nil, nil
}

func InitializeServiceForTest() (*Service, error) {
	wire.Build(
		categoryrepositorylocal.NewCategoryRepositoryLocal, // Provider para CategoryRepository
		quantityusecases.InitializeServiceForTest,          // Provider para QuantityRepository
		sizeusecases.InitializeServiceForTest,
		NewService, // Construtor do serviço principal
	)
	return nil, nil
}
