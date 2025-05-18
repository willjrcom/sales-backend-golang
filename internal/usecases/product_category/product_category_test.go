package productcategoryusecases

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"

	categoryrepolocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/category"
	quantityrepolocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/quantity"
	sizerepolocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/size"

	quantityusecases "github.com/willjrcom/sales-backend-go/internal/usecases/quantity"
	sizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/size"
)

var (
	service *Service
	ctx     context.Context
)

func TestMain(m *testing.M) {
	ctx = context.Background()

	// Shared repositories for test isolation
	catRepo := categoryrepolocal.NewCategoryRepositoryLocal()
	qtyRepo := quantityrepolocal.NewQuantityRepositoryLocal()
	szRepo := sizerepolocal.NewSizeRepositoryLocal()

	// Use-case services with shared category repository
	qtySvc := quantityusecases.NewService(qtyRepo, catRepo)
	szSvc := sizeusecases.NewService(szRepo, catRepo)
	service = NewService(catRepo, qtySvc, szSvc)

	m.Run()
}

func TestCreateCategory(t *testing.T) {
	dto := &productcategorydto.CategoryCreateDTO{Name: "test-category"}
	catID, err := service.CreateCategory(ctx, dto)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, catID)

	dtoID := entitydto.NewIdRequest(catID)
	cat, err := service.GetCategoryById(ctx, dtoID)
	assert.NoError(t, err)
	assert.Equal(t, dto.Name, cat.Name)
}