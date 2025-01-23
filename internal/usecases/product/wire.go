//go:build wireinject
// +build wireinject

package productusecases

import (
	"github.com/google/wire"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	productrepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/product"
	productcategoryrepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/product_category"
	productrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product"
	categoryrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category"
	s3service "github.com/willjrcom/sales-backend-go/internal/infra/service/s3"
)

func InitializeService() (*Service, error) {
	wire.Build(
		database.NewPostgreSQLConnection,                      // Provider para *bun.DB
		productrepositorybun.NewProductRepositoryBun,          // Provider para ProductRepository
		categoryrepositorybun.NewProductCategoryRepositoryBun, // Provider para CategoryRepository
		s3service.NewS3Client,                                 // Provider para S3Client
		NewService,                                            // Construtor do serviço principal
	)
	return nil, nil
}

func InitializeServiceForTest() (*Service, error) {
	wire.Build(
		productrepositorylocal.NewProductRepositoryLocal,          // Provider para ProductRepository
		productcategoryrepositorylocal.NewCategoryRepositoryLocal, // Provider para CategoryRepository
		s3service.NewS3Client,                                     // Provider para S3Client
		NewService,                                                // Construtor do serviço principal
	)
	return nil, nil
}
