package productusecases

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	filterdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/filter"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
	productrepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/product"
)

var service *Service

func TestMain(m *testing.M) {
	repo := productrepositorylocal.NewProductRepositoryLocal()
	service = NewService(repo)

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestRegisterProduct(t *testing.T) {
	dto := &productdto.CreateProductInput{
		Code:     "123",
		Name:     "Test Product",
		Cost:     100,
		Price:    100,
		Category: "Test Category",
	}

	idProduct, err := service.RegisterProduct(dto)

	dtoId := entitydto.NewIdRequest(idProduct)
	product, err := service.GetProductById(dtoId)

	assert.Nil(t, err)
	assert.NotContains(t, idProduct, uuid.Nil)
	assert.Equal(t, product.Name, "Test Product")
	assert.Equal(t, product.ID, idProduct)
}

func TestRegisterProductError(t *testing.T) {
	// Teste 1 - No Code
	dto := &productdto.CreateProductInput{
		Name:     "Test Product",
		Cost:     90,
		Price:    100,
		Category: "Test Category",
	}

	_, err := service.RegisterProduct(dto)
	assert.NotNil(t, err)
	assert.EqualError(t, err, productdto.ErrCodeRequired.Error())

	// Test 2 - No Name
	dto = &productdto.CreateProductInput{
		Code:     "CODE",
		Cost:     90,
		Price:    100,
		Category: "Test Category",
	}

	_, err = service.RegisterProduct(dto)
	assert.NotNil(t, err)
	assert.EqualError(t, err, productdto.ErrNameRequired.Error())

	// Test 3 - Price greater than cost
	dto = &productdto.CreateProductInput{
		Code:     "CODE",
		Name:     "Test Product",
		Cost:     150,
		Price:    100,
		Category: "Test Category",
	}

	_, err = service.RegisterProduct(dto)
	assert.NotNil(t, err)
	assert.EqualError(t, err, productdto.ErrCostGreaterThanPrice.Error())

	// Test 4 - No category
	dto = &productdto.CreateProductInput{
		Code:  "CODE",
		Name:  "Test Product",
		Cost:  90,
		Price: 100,
	}

	_, err = service.RegisterProduct(dto)
	assert.NotNil(t, err)
	assert.EqualError(t, err, productdto.ErrCategoryRequired.Error())
}

func TestUpdateProduct(t *testing.T) {
	products, err := service.GetAllProduct(&filterdto.Filter{})

	assert.Equal(t, len(products), 1)
	idProduct := products[0].ID

	dto := &productdto.UpdateProductInput{}
	dtoId := entitydto.NewIdRequest(idProduct)

	jsonTest1 := []byte(`{"name": "new Product"}`)
	jsonTest2 := []byte(`{"cost": 150}`)
	jsonTest3 := []byte(`{"category": ""}`)

	// Test 1 - New name
	assert.Nil(t, json.Unmarshal(jsonTest1, &dto))
	assert.Equal(t, "new Product", (*dto.Name))

	err = service.UpdateProduct(dtoId, dto)
	assert.Nil(t, err)

	// Test 2 - Cost greater than Price
	assert.Nil(t, json.Unmarshal(jsonTest2, &dto))

	err = service.UpdateProduct(dtoId, dto)
	assert.EqualError(t, err, productdto.ErrCostGreaterThanPrice.Error())
	*dto.Cost = float64(90.0)

	// Test 3 - No category
	assert.Nil(t, json.Unmarshal(jsonTest3, &dto))

	err = service.UpdateProduct(dtoId, dto)
	assert.EqualError(t, err, productdto.ErrCategoryRequired.Error())
}

func TestGetAll(t *testing.T) {
	products, err := service.GetAllProduct(&filterdto.Filter{})

	assert.Nil(t, err)
	assert.Equal(t, 1, len(products))
}

func TestGetProductById(t *testing.T) {
	products, err := service.GetAllProduct(&filterdto.Filter{})
	assert.Equal(t, len(products), 1)
	idProduct := products[0].ID

	dtoId := entitydto.NewIdRequest(idProduct)
	product, err := service.GetProductById(dtoId)

	assert.Nil(t, err)
	assert.Equal(t, "new Product", product.Name)
}
