package orderusecases

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	filterdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/filter"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
	orderrepositorylocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/order"
)

var service *Service

func TestMain(m *testing.M) {
	repo := orderrepositorylocal.NewOrderRepositoryLocal()
	service = NewService(repo)

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestRegisterOrder(t *testing.T) {
	dto := &orderdto.CreateOrderInput{
		Name: "Test Order",
	}

	idOrder, err := service.CreateOrder(dto)

	dtoId := entitydto.NewIdRequest(idOrder)
	Order, err := service.GetOrderById(dtoId)

	assert.Nil(t, err)
	assert.NotContains(t, idOrder, uuid.Nil)
	assert.Equal(t, Order.Name, "Test Order")
	assert.Equal(t, Order.ID, idOrder)
}

func TestRegisterOrderError(t *testing.T) {
	// Teste 1 - No Name
	dto := &orderdto.CreateOrderInput{}

	_, err := service.CreateOrder(dto)
	assert.NotNil(t, err)
	assert.EqualError(t, err, orderdto.ErrNameRequired.Error())

}

func TestUpdateOrder(t *testing.T) {

	Orders, err := service.GetAllOrder(&filterdto.Filter{})

	assert.Equal(t, len(Orders), 1)
	idOrder := Orders[0].ID

	dto := &orderdto.UpdateObservationOrder{Observation: "New observation"}
	dtoId := entitydto.NewIdRequest(idOrder)

	// Test 1 - New observation
	err = service.UpdateOrderObservation(dtoId, dto)
	assert.Nil(t, err)
}

func TestGetAll(t *testing.T) {
	Orders, err := service.GetAllOrder(&filterdto.Filter{})

	assert.Nil(t, err)
	assert.Equal(t, 1, len(Orders))
}

func TestGetOrderById(t *testing.T) {
	Orders, err := service.GetAllOrder(&filterdto.Filter{})
	assert.Equal(t, len(Orders), 1)
	idOrder := Orders[0].ID

	dtoId := entitydto.NewIdRequest(idOrder)
	Order, err := service.GetOrderById(dtoId)

	assert.Nil(t, err)
	assert.Equal(t, "Test Order", Order.Name)
}
