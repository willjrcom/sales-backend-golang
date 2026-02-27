package pos

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

// Test FormatOrder starts and ends with correct control codes.
func Test_FormatOrder(t *testing.T) {
	now := time.Now().UTC()
	decimalValue := decimal.NewFromFloat(10.00)

	o := &orderentity.Order{
		Entity: entity.Entity{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		OrderCommonAttributes: orderentity.OrderCommonAttributes{
			OrderNumber: 1,
			OrderType: orderentity.OrderType{
				Delivery: &orderentity.OrderDelivery{
					Entity: entity.Entity{
						ID: uuid.New(),
					},
					DeliveryTimeLogs: orderentity.DeliveryTimeLogs{
						PendingAt: &now,
						ShippedAt: &now,
					},
					OrderDeliveryCommonAttributes: orderentity.OrderDeliveryCommonAttributes{
						Client: &cliententity.Client{
							Person: personentity.Person{
								PersonCommonAttributes: personentity.PersonCommonAttributes{
									Name: "Client Will",
								},
							},
						},
						Address: &addressentity.Address{
							AddressCommonAttributes: addressentity.AddressCommonAttributes{
								Street:       "Street",
								Number:       "123",
								Complement:   "Complement",
								Reference:    "Reference",
								Neighborhood: "Neighborhood",
								City:         "City",
								UF:           "UF",
								Cep:          "00000-000",
							},
						},
						Driver: &orderentity.DeliveryDriver{
							DeliveryDriverCommonAttributes: orderentity.DeliveryDriverCommonAttributes{
								Employee: &employeeentity.Employee{
									User: &companyentity.User{
										UserCommonAttributes: companyentity.UserCommonAttributes{
											Person: personentity.Person{
												PersonCommonAttributes: personentity.PersonCommonAttributes{
													Name: "Motoboy Will",
												},
											},
										},
									},
								},
							},
						},
						DeliveryTax:   &decimalValue,
						Change:        decimalValue,
						PaymentMethod: orderentity.Dinheiro,
					},
				},
			},
			OrderDetail: orderentity.OrderDetail{
				SubTotal:      decimal.NewFromFloat(150.45),
				TotalPaid:     decimal.NewFromFloat(120.45),
				TotalChange:   decimal.NewFromFloat(30.45),
				QuantityItems: 2.0,
				Observation:   "Order note",
				Attendant: &employeeentity.Employee{
					User: &companyentity.User{
						UserCommonAttributes: companyentity.UserCommonAttributes{
							Person: personentity.Person{
								PersonCommonAttributes: personentity.PersonCommonAttributes{
									Name: "John Doe",
								},
							},
						},
					},
				},
			},
			GroupItems: []orderentity.GroupItem{
				{
					GroupCommonAttributes: orderentity.GroupCommonAttributes{
						Items: []orderentity.Item{
							{
								ItemCommonAttributes: orderentity.ItemCommonAttributes{
									Name:        "Pizza mussarela",
									Observation: "bem quente",
									Size:        "G",
									TotalPrice:  decimal.NewFromFloat(30.00),
									RemovedItems: []string{
										"Tomate",
									},
									Quantity: 0.5,
									AdditionalItems: []orderentity.Item{
										{
											ItemCommonAttributes: orderentity.ItemCommonAttributes{
												Name:       "Bacon",
												TotalPrice: decimal.NewFromFloat(10.00),
												Quantity:   1.0,
											},
										},
										{
											ItemCommonAttributes: orderentity.ItemCommonAttributes{
												Name:       "Cebola",
												TotalPrice: decimal.NewFromFloat(10.00),
												Quantity:   1.0,
											},
										},
									},
								},
							},
							{
								ItemCommonAttributes: orderentity.ItemCommonAttributes{
									Name:       "Pizza calabresa",
									TotalPrice: decimal.NewFromFloat(30.00),
									Quantity:   0.5,
									Size:       "G",
									RemovedItems: []string{
										"Cebola",
									},
								},
							},
						},
						GroupDetails: orderentity.GroupDetails{
							Size:        "G",
							TotalPrice:  decimal.NewFromFloat(100.00),
							Quantity:    1.0,
							Observation: "sem molho",
							Category: &productentity.ProductCategory{
								ProductCategoryCommonAttributes: productentity.ProductCategoryCommonAttributes{
									Name: "Pizza",
								},
							},
							ComplementItem: &orderentity.Item{
								ItemCommonAttributes: orderentity.ItemCommonAttributes{
									Name:        "Borda Recheada",
									TotalPrice:  decimal.NewFromFloat(10.00),
									Observation: "bastante catupiry",
									Size:        "G",
									Quantity:    1.0,
								},
							},
						},
					},
					GroupItemTimeLogs: orderentity.GroupItemTimeLogs{
						StartAt:     &now,
						PendingAt:   &now,
						StartedAt:   &now,
						ReadyAt:     &now,
						CancelledAt: &now,
					},
				},
			},
			Payments: []orderentity.PaymentOrder{
				{
					PaymentCommonAttributes: orderentity.PaymentCommonAttributes{
						TotalPaid: decimal.NewFromFloat(30.45),
						Method:    orderentity.Dinheiro,
					},
				},
				{
					PaymentCommonAttributes: orderentity.PaymentCommonAttributes{
						TotalPaid: decimal.NewFromFloat(20),
						Method:    orderentity.Alelo,
					},
				},
			},
		},
	}

	out, err := FormatOrder(o, nil)
	assert.NoError(t, err)
	if err := os.WriteFile("printer_order.txt", out, 0644); err != nil {
		t.Fatalf("failed to write printer buffer to file: %v", err)
	}
	s := string(out)
	// Deve iniciar com código ESC inicial e conter o cabeçalho de entrega
	assert.True(t, strings.HasPrefix(s, escInit))
	assert.Contains(t, s, fmt.Sprintf("PEDIDO DE ENTREGA %d", o.OrderNumber))
	// Deve conter seção de entrega e informações principais
	assert.Contains(t, s, "ENTREGA")
	assert.Contains(t, s, "Despachado:")
	assert.Contains(t, s, "Cliente:")
	assert.Contains(t, s, "Street, 123")
	assert.Contains(t, s, "Complemento:")
	assert.Contains(t, s, "Ref:")
	assert.Contains(t, s, "Bairro:")
	assert.Contains(t, s, "Cidade:")
	assert.Contains(t, s, "CEP:")
	assert.Contains(t, s, "Entregador:")
	assert.Contains(t, s, "Taxa entrega:")
	assert.Contains(t, s, "Troco:")
	assert.Contains(t, s, "Forma de pagamento:")
	// Itens do pedido e total
	assert.Contains(t, s, "Pizza mussarela")
	assert.Contains(t, s, "Borda Recheada")
	assert.Contains(t, s, "TOTAL:")
	// Deve terminar com código de corte
	assert.True(t, strings.HasSuffix(s, escCut))
}
