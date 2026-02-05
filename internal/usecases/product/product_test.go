package productusecases

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size"

	productrepolocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/product"
	categoryrepolocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/product_category"
	sizerepolocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/size"

	productcategoryusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category"
	sizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/size"

	s3service "github.com/willjrcom/sales-backend-go/internal/infra/service/s3"
)

var (
	productService         *Service
	sizeService            *sizeusecases.Service
	productCategoryService *productcategoryusecases.Service
	ctx                    context.Context
)

// ptrUUID returns a pointer to the given UUID.
func ptrUUID(u uuid.UUID) *uuid.UUID {
	return &u
}

func TestMain(m *testing.M) {
	ctx = context.Background()

	// Shared in-memory repositories for test isolation
	catRepo := categoryrepolocal.NewCategoryRepositoryLocal()
	szRepo := sizerepolocal.NewSizeRepositoryLocal()
	prdRepo := productrepolocal.NewProductRepositoryLocal()

	// Use-case services with shared category repository
	szSvc := sizeusecases.NewService(szRepo, catRepo)
	productCategoryService = productcategoryusecases.NewService(catRepo, szSvc)
	sizeService = sizeusecases.NewService(szRepo, catRepo)
	productService = NewService(prdRepo, catRepo, s3service.NewS3Client())

	os.Exit(m.Run())
}

func TestCreateProduct(t *testing.T) {
	dtoCategory := &productcategorydto.CategoryCreateDTO{Name: "pizza"}
	categoryId, err := productCategoryService.CreateCategory(ctx, dtoCategory)
	assert.Nil(t, err)
	assert.NotNil(t, categoryId)

	dtoSize := &sizedto.SizeCreateDTO{
		Name:       "P",
		CategoryID: categoryId,
	}
	sizeId, err := sizeService.CreateSize(ctx, dtoSize)
	assert.Nil(t, err)
	assert.NotNil(t, sizeId)

	dto := &productcategorydto.ProductCreateDTO{
		CategoryID: &categoryId,
		SizeID:     &sizeId,
		Code:       "test",
		Name:       "Test Product",
	}

	productId, err := productService.CreateProduct(ctx, dto)
	assert.Nil(t, err)

	dtoId := entitydto.NewIdRequest(productId)
	product, err := productService.GetProductById(ctx, dtoId)

	assert.Nil(t, err)
	assert.NotContains(t, productId, uuid.Nil)
	assert.Equal(t, product.Name, "Test Product")
	assert.Equal(t, product.ID, productId)
	assert.NotNil(t, product.Flavors)
	assert.Len(t, product.Flavors, 0)
}

func TestCreateProductError(t *testing.T) {
	// Test 1 - No Code
	dto := &productcategorydto.ProductCreateDTO{}
	_, err := productService.CreateProduct(ctx, dto)
	assert.EqualError(t, err, productcategorydto.ErrCodeRequired.Error())

	// Test 2 - No Name
	dto = &productcategorydto.ProductCreateDTO{
		Code:       "code",
		CategoryID: ptrUUID(uuid.Nil),
		SizeID:     ptrUUID(uuid.Nil),
	}
	_, err = productService.CreateProduct(ctx, dto)
	assert.EqualError(t, err, productcategorydto.ErrNameRequired.Error())

	// Test 3 - Price less than Cost
	dto = &productcategorydto.ProductCreateDTO{
		Code:       "code",
		Name:       "name",
		Price:      decimal.NewFromInt(1),
		Cost:       decimal.NewFromInt(2),
		CategoryID: ptrUUID(uuid.Nil),
		SizeID:     ptrUUID(uuid.Nil),
	}
	_, err = productService.CreateProduct(ctx, dto)
	assert.EqualError(t, err, productcategorydto.ErrCostGreaterThanPrice.Error())

	// Test 4 - No Category
	dto = &productcategorydto.ProductCreateDTO{
		Code:   "code",
		Name:   "name",
		Price:  decimal.NewFromInt(2),
		Cost:   decimal.NewFromInt(1),
		SizeID: ptrUUID(uuid.Nil),
	}
	_, err = productService.CreateProduct(ctx, dto)
	assert.EqualError(t, err, productcategorydto.ErrCategoryRequired.Error())

	// Test 5 - No Size
	dto = &productcategorydto.ProductCreateDTO{
		Code:       "code",
		Name:       "name",
		Price:      decimal.NewFromInt(2),
		Cost:       decimal.NewFromInt(1),
		CategoryID: ptrUUID(uuid.Nil),
	}
	_, err = productService.CreateProduct(ctx, dto)
	assert.EqualError(t, err, productcategorydto.ErrSizeRequired.Error())
}

func TestUpdateProduct(t *testing.T) {
	products, _, err := productService.GetAllProducts(ctx, 0, 100, true, "")

	assert.Nil(t, err)
	assert.NotNil(t, products)
	assert.Equal(t, len(products), 1)

	idProduct := products[0].ID

	dto := &productcategorydto.ProductUpdateDTO{}
	dtoId := entitydto.NewIdRequest(idProduct)

	jsonTest1 := []byte(`{"name": "new Product"}`)
	jsonTest2 := []byte(`{"cost": 150}`)

	// Test 1 - New name
	assert.Nil(t, json.Unmarshal(jsonTest1, &dto))
	assert.Equal(t, "new Product", (*dto.Name))

	err = productService.UpdateProduct(ctx, dtoId, dto)
	assert.Nil(t, err)

	// Test 2 - Cost greater than Price
	assert.Nil(t, json.Unmarshal(jsonTest2, &dto))

	err = productService.UpdateProduct(ctx, dtoId, dto)
	assert.EqualError(t, err, productcategorydto.ErrCostGreaterThanPrice.Error())
	*dto.Cost = decimal.NewFromInt(90)

	// Test 3 - Update flavors
	dto = &productcategorydto.ProductUpdateDTO{
		Flavors: []string{"mussarela", "calabresa"},
	}
	err = productService.UpdateProduct(ctx, dtoId, dto)
	assert.Nil(t, err)

	product, err := productService.GetProductById(ctx, dtoId)
	assert.Nil(t, err)
	assert.Equal(t, []string{"mussarela", "calabresa"}, product.Flavors)
}

func TestGetAll(t *testing.T) {
	products, _, err := productService.GetAllProducts(ctx, 0, 100, true, "")

	assert.Nil(t, err)
	assert.Equal(t, 1, len(products))
}

func TestGetProductById(t *testing.T) {
	products, _, _ := productService.GetAllProducts(ctx, 0, 100, true, "")
	assert.Equal(t, len(products), 1)
	idProduct := products[0].ID

	dtoId := entitydto.NewIdRequest(idProduct)
	product, err := productService.GetProductById(ctx, dtoId)

	assert.Nil(t, err)
	assert.Equal(t, "new Product", product.Name)
}
