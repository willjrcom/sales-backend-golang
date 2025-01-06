package productusecases

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size"
	productrepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/product"
	categoryrepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/product_category"
	sizerepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/size"
	s3service "github.com/willjrcom/sales-backend-go/internal/infra/service/s3"
	categoryusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category"
	categorysizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/size"
)

var (
	productService         *Service
	sizeService            *categorysizeusecases.Service
	productCategoryService *categoryusecases.Service
	ctx                    context.Context
)

func TestMain(m *testing.M) {
	ctx = context.Background()

	// Repository
	productrepositorylocal.NewProductRepositoryLocal()
	rc := categoryrepositorylocal.NewCategoryRepositoryLocal()
	rs := sizerepositorylocal.NewSizeRepositoryLocal()
	s3service.NewS3Client()

	// Service
	productService, _ = InitializeService()
	sizeService = categorysizeusecases.NewService(rs)

	sizeService.AddDependencies(rc)

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestCreateProduct(t *testing.T) {
	dtoCategory := &productcategorydto.CreateCategoryInput{ProductCategoryCommonAttributes: productentity.ProductCategoryCommonAttributes{Name: "pizza"}}
	categoryId, err := productCategoryService.CreateCategory(ctx, dtoCategory)
	assert.Nil(t, err)
	assert.NotNil(t, categoryId)

	dtoSize := &sizedto.SizeCreateDTO{Name: "P"}
	sizeId, err := sizeService.CreateSize(ctx, dtoSize)
	assert.Nil(t, err)
	assert.NotNil(t, sizeId)

	dto := &productdto.CreateProductInput{
		PatchProduct: productentity.PatchProduct{},
	}

	productId, err := productService.CreateProduct(ctx, dto)
	assert.Nil(t, err)

	dtoId := entitydto.NewIdRequest(productId)
	product, err := productService.GetProductById(ctx, dtoId)

	assert.Nil(t, err)
	assert.NotContains(t, productId, uuid.Nil)
	assert.Equal(t, product.Name, "Test Product")
	assert.Equal(t, product.ID, productId)
}

func TestCreateProductError(t *testing.T) {
	// Teste 1 - No Code
	dto := &productdto.CreateProductInput{
		PatchProduct: productentity.PatchProduct{},
	}

	_, err := productService.CreateProduct(ctx, dto)
	assert.NotNil(t, err)
	assert.EqualError(t, err, productdto.ErrCodeRequired.Error())

	// Test 2 - No Name
	dto = &productdto.CreateProductInput{
		PatchProduct: productentity.PatchProduct{},
	}

	_, err = productService.CreateProduct(ctx, dto)
	assert.NotNil(t, err)
	assert.EqualError(t, err, productdto.ErrNameRequired.Error())

	// Test 3 - Price greater than cost
	dto = &productdto.CreateProductInput{
		PatchProduct: productentity.PatchProduct{},
	}

	_, err = productService.CreateProduct(ctx, dto)
	assert.NotNil(t, err)
	assert.EqualError(t, err, productdto.ErrCostGreaterThanPrice.Error())

	// Test 4 - No category
	dto = &productdto.CreateProductInput{
		PatchProduct: productentity.PatchProduct{},
	}

	_, err = productService.CreateProduct(ctx, dto)
	assert.NotNil(t, err)
	assert.EqualError(t, err, productdto.ErrCategoryRequired.Error())
}

func TestUpdateProduct(t *testing.T) {
	products, err := productService.GetAllProducts(ctx)

	assert.Nil(t, err)
	assert.NotNil(t, products)
	assert.Equal(t, len(products), 1)

	idProduct := products[0].ID

	dto := &productdto.UpdateProductInput{}
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
	assert.EqualError(t, err, productdto.ErrCostGreaterThanPrice.Error())
	*dto.Cost = float64(90.0)
}

func TestGetAll(t *testing.T) {
	products, err := productService.GetAllProducts(ctx)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(products))
}

func TestGetProductById(t *testing.T) {
	products, _ := productService.GetAllProducts(ctx)
	assert.Equal(t, len(products), 1)
	idProduct := products[0].ID

	dtoId := entitydto.NewIdRequest(idProduct)
	product, err := productService.GetProductById(ctx, dtoId)

	assert.Nil(t, err)
	assert.Equal(t, "new Product", product.Name)
}
