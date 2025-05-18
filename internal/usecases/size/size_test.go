package sizeusecases_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size"

	categoryrepolocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/category"
	quantityrepolocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/quantity"
	sizerepolocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/size"

	productcategoryusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category"
	quantityusecases "github.com/willjrcom/sales-backend-go/internal/usecases/quantity"
	sizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/size"
)

var (
	service                *sizeusecases.Service
	productCategoryService *productcategoryusecases.Service
	quantityService        *quantityusecases.Service
	szctx                  context.Context
)

func TestMain(m *testing.M) {
	szctx = context.Background()

	// Shared in-memory repositories
	catRepo := categoryrepolocal.NewCategoryRepositoryLocal()
	qtyRepo := quantityrepolocal.NewQuantityRepositoryLocal()
	szRepo := sizerepolocal.NewSizeRepositoryLocal()

	// Initialize dependent services
	quantityService = quantityusecases.NewService(qtyRepo, catRepo)
	service = sizeusecases.NewService(szRepo, catRepo)
	productCategoryService = productcategoryusecases.NewService(catRepo, quantityService, service)

	os.Exit(m.Run())
}

func TestCreateSize(t *testing.T) {
	// First create a category to own sizes
	dtoCat := &productcategorydto.CategoryCreateDTO{Name: "pizza"}
	catID, err := productCategoryService.CreateCategory(szctx, dtoCat)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, catID)

	// Create a size under the category (distinct from default sizes)
	dto := &sizedto.SizeCreateDTO{Name: "X", CategoryID: catID}
	szID, err := service.CreateSize(szctx, dto)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, szID)

	dtoId := entitydto.NewIdRequest(szID)
	sz, err := service.GetSizeById(szctx, dtoId)
	assert.NoError(t, err)
	assert.Equal(t, dto.Name, sz.Name)
}
