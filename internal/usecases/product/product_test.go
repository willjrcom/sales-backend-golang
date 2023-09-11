package productusecases

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	filterdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/filter"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
	categoryrepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/category-product"
	productrepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/product"
	categoryproductusecases "github.com/willjrcom/sales-backend-go/internal/usecases/category_product"
)

var (
	productService         *Service
	categoryProductService *categoryproductusecases.Service
	ctx                    context.Context
)

func TestMain(m *testing.M) {
	ctx = context.Background()

	// Repository
	rProduct := productrepositorylocal.NewProductRepositoryLocal()
	rCategoryProduct := categoryrepositorylocal.NewCategoryProductRepositoryLocal()

	// Service
	productService = NewService(rProduct, rCategoryProduct)

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestRegisterProduct(t *testing.T) {
	dtoCategory := &productdto.RegisterCategoryProductInput{Name: "pizza", Sizes: []string{"P", "M", "G"}}

	categoryId, err := categoryProductService.RegisterCategoryProduct(ctx, dtoCategory)

	assert.Nil(t, err)
	assert.NotNil(t, categoryId)

	dto := &productdto.RegisterProductInput{
		Code:       "123",
		Name:       "Test Product",
		Cost:       100,
		Price:      100,
		CategoryID: categoryId,
	}

	productId, err := productService.RegisterProduct(ctx, dto)
	assert.Nil(t, err)

	dtoId := entitydto.NewIdRequest(productId)
	product, err := productService.GetProductById(ctx, dtoId)

	assert.Nil(t, err)
	assert.NotContains(t, productId, uuid.Nil)
	assert.Equal(t, product.Name, "Test Product")
	assert.Equal(t, product.ID, productId)
}

func TestRegisterProductError(t *testing.T) {
	// Teste 1 - No Code
	dto := &productdto.RegisterProductInput{
		Name:       "Test Product",
		Cost:       90,
		Price:      100,
		CategoryID: uuid.New(),
	}

	_, err := productService.RegisterProduct(ctx, dto)
	assert.NotNil(t, err)
	assert.EqualError(t, err, productdto.ErrCodeRequired.Error())

	// Test 2 - No Name
	dto = &productdto.RegisterProductInput{
		Code:       "CODE",
		Cost:       90,
		Price:      100,
		CategoryID: uuid.New(),
	}

	_, err = productService.RegisterProduct(ctx, dto)
	assert.NotNil(t, err)
	assert.EqualError(t, err, productdto.ErrNameRequired.Error())

	// Test 3 - Price greater than cost
	dto = &productdto.RegisterProductInput{
		Code:       "CODE",
		Name:       "Test Product",
		Cost:       150,
		Price:      100,
		CategoryID: uuid.New(),
	}

	_, err = productService.RegisterProduct(ctx, dto)
	assert.NotNil(t, err)
	assert.EqualError(t, err, productdto.ErrCostGreaterThanPrice.Error())

	// Test 4 - No category
	dto = &productdto.RegisterProductInput{
		Code:  "CODE",
		Name:  "Test Product",
		Cost:  90,
		Price: 100,
	}

	_, err = productService.RegisterProduct(ctx, dto)
	assert.NotNil(t, err)
	assert.EqualError(t, err, productdto.ErrCategoryRequired.Error())
}

func TestUpdateProduct(t *testing.T) {
	products, err := productService.GetAllProductsByCategory(ctx, &filterdto.Category{Category: "teste"})

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
	products, err := productService.GetAllProductsByCategory(ctx, &filterdto.Category{Category: "teste"})

	assert.Nil(t, err)
	assert.Equal(t, 1, len(products))
}

func TestGetProductById(t *testing.T) {
	products, _ := productService.GetAllProductsByCategory(ctx, &filterdto.Category{Category: "teste"})
	assert.Equal(t, len(products), 1)
	idProduct := products[0].ID

	dtoId := entitydto.NewIdRequest(idProduct)
	product, err := productService.GetProductById(ctx, dtoId)

	assert.Nil(t, err)
	assert.Equal(t, "new Product", product.Name)
}
