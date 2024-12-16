//go:build wireinject
// +build wireinject

package productcategoryproductusecases

import (
	"github.com/google/wire"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	productcategoryrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category"
	productcategoryproductrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/product_category_product"
	s3service "github.com/willjrcom/sales-backend-go/internal/infra/service/s3"
)

func InitializeService() (*Service, error) {
	wire.Build(
		database.NewPostgreSQLConnection,                             // Provider para *bun.DB
		productcategoryproductrepositorybun.NewProductRepositoryBun,  // Provider para ProductRepository
		productcategoryrepositorybun.NewProductCategoryRepositoryBun, // Provider para CategoryRepository
		s3service.NewS3Client,                                        // Provider para S3Client
		NewService,                                                   // Construtor do serviço principal
	)
	return nil, nil
}
