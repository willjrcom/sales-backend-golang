package pos

import (
	"bytes"
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

// Test printAdditionalItem prints a single additional item line correctly.
func Test_printAdditionalItem(t *testing.T) {
	buf := &bytes.Buffer{}
	add := orderentity.NewItem("Extra Cheese", decimal.NewFromFloat(2.50), 1, "", uuid.New(), uuid.New())
	printAdditionalItem(buf, add)
	expected := fmt.Sprintf("   + %-17s %7.2f%s", add.Name, d2f(add.TotalPrice), newline)
	assert.Equal(t, expected, buf.String())
}

// Test printItem prints main item line and handles observation and removed items.
func Test_printItem_WithExtras(t *testing.T) {
	buf := &bytes.Buffer{}
	// Main item
	item := orderentity.NewItem("Pizza", decimal.NewFromFloat(5.00), 2, "", uuid.New(), uuid.New())
	// Additional item
	extra := orderentity.NewItem("Bacon", decimal.NewFromFloat(1.50), 1, "", uuid.New(), uuid.New())
	item.AdditionalItems = []orderentity.Item{*extra}
	// Observation and removed items
	item.Observation = "No onions"
	item.RemovedItems = []string{"Onion"}
	printItem(buf, item)
	// Build expected output
	expected := fmt.Sprintf("%2.0f x %-20s %7.2f%s", item.Quantity, item.Name, d2f(item.TotalPrice), newline)
	expected += fmt.Sprintf("   + %-17s %7.2f%s", extra.Name, d2f(extra.TotalPrice), newline)
	expected += fmt.Sprintf("   Obs: %s%s", item.Observation, newline)
	expected += fmt.Sprintf("   Removido: %s%s", item.RemovedItems[0], newline)
	assert.Equal(t, expected, buf.String())
}

// Test printItem without extras prints only the main line.
func Test_printItem_NoExtras(t *testing.T) {
	buf := &bytes.Buffer{}
	item := orderentity.NewItem("Sandwich", decimal.NewFromFloat(3.75), 1, "", uuid.New(), uuid.New())
	printItem(buf, item)
	expected := fmt.Sprintf("%2.0f x %-20s %7.2f%s", item.Quantity, item.Name, d2f(item.TotalPrice), newline)
	assert.Equal(t, expected, buf.String())
}

// Test printGroupItem prints header, quantity, observation, items, and subtotal.
func Test_printGroupItem(t *testing.T) {
	buf := &bytes.Buffer{}
	// Setup category
	cat := productentity.NewProductCategory(productentity.ProductCategoryCommonAttributes{Name: "Categoria"})
	// Create group and add item
	group := orderentity.NewGroupItem(orderentity.GroupCommonAttributes{GroupDetails: orderentity.GroupDetails{
		Category: cat, Size: "M",
	}})
	item := orderentity.NewItem("Produto", decimal.NewFromFloat(10.00), 2, "", uuid.New(), uuid.New())
	group.Items = []orderentity.Item{*item}
	// Calculate totals and set observation
	group.CalculateTotalPrice()
	group.Observation = "Group note"
	printGroupItem(buf, group)
	// Build expected output
	header := fmt.Sprintf("%-40s%s", "Categoria - M", newline)
	expected := escAlignLeft + escBoldOn + header + escBoldOff
	expected += fmt.Sprintf("Qtd:%36.0f%s", group.Quantity, newline)
	expected += fmt.Sprintf("Obs: %s%s", group.Observation, newline)
	expected += fmt.Sprintf("%2.0f x %-20s %7.2f%s", item.Quantity, item.Name, d2f(item.TotalPrice), newline)
	expected += fmt.Sprintf("Subtotal:%31.2f%s", d2f(group.TotalPrice), newline)
	assert.Equal(t, expected, buf.String())
}

// Test printGroupItemsSection prints multiple groups and a separator line.
func Test_printGroupItemsSection(t *testing.T) {
	// Single group scenario matches FormatGroupItem
	cat := productentity.NewProductCategory(productentity.ProductCategoryCommonAttributes{Name: "Cat"})
	group := orderentity.NewGroupItem(orderentity.GroupCommonAttributes{GroupDetails: orderentity.GroupDetails{
		Category: cat, Size: "S",
	}})
	item := orderentity.NewItem("Item1", decimal.NewFromFloat(4.00), 1, "", uuid.New(), uuid.New())
	group.Items = []orderentity.Item{*item}
	group.CalculateTotalPrice()
	// Expected from FormatGroupItem
	expectedBytes, _ := FormatGroupItem(group)
	buf := &bytes.Buffer{}
	printGroupItemsSection(buf, []orderentity.GroupItem{*group})
	assert.Equal(t, string(expectedBytes), buf.String())
}

// Test formatHeader prints order number and separator.
func Test_formatHeader(t *testing.T) {
	buf := &bytes.Buffer{}
	o := &orderentity.Order{OrderCommonAttributes: orderentity.OrderCommonAttributes{OrderNumber: 42}}
	formatHeader(buf, o)
	expected := escAlignCenter + escBoldOn + fmt.Sprintf("PEDIDO %d%s", o.OrderNumber, newline) + escBoldOff + escAlignLeft + strings.Repeat("-", 40) + newline
	assert.Equal(t, expected, buf.String())
}

// Test formatFooter prints total payable.
func Test_formatFooter(t *testing.T) {
	buf := &bytes.Buffer{}
	o := &orderentity.Order{OrderCommonAttributes: orderentity.OrderCommonAttributes{}}
	o.TotalPayable = decimal.NewFromFloat(123.45)
	formatFooter(buf, o)
	expected := fmt.Sprintf("TOTAL:%31.2f%s", d2f(o.TotalPayable), newline)
	assert.Equal(t, expected, buf.String())
}

// Test FormatItem delegates to printItem.
func Test_FormatItem(t *testing.T) {
	item := orderentity.NewItem("Burger", decimal.NewFromFloat(6.00), 1, "", uuid.New(), uuid.New())
	item.Observation = "Extra sauce"
	out, err := FormatItem(item)
	assert.NoError(t, err)
	buf := &bytes.Buffer{}
	printItem(buf, item)
	assert.Equal(t, buf.String(), string(out))
}

// Test FormatGroupItem delegates to printGroupItem and adds separator.
func Test_FormatGroupItem(t *testing.T) {
	cat := productentity.NewProductCategory(productentity.ProductCategoryCommonAttributes{Name: "Cat"})
	group := orderentity.NewGroupItem(orderentity.GroupCommonAttributes{GroupDetails: orderentity.GroupDetails{
		Category: cat, Size: "L",
	}})
	item := orderentity.NewItem("Item2", decimal.NewFromFloat(7.50), 3, "", uuid.New(), uuid.New())
	group.Items = []orderentity.Item{*item}
	group.CalculateTotalPrice()
	out, err := FormatGroupItem(group)
	assert.NoError(t, err)
	buf := &bytes.Buffer{}
	printGroupItem(buf, group)
	buf.WriteString(strings.Repeat("-", 40) + newline)
	assert.Equal(t, buf.String(), string(out))
}

// Test FormatOrder starts and ends with correct control codes.
func Test_FormatOrder(t *testing.T) {
	now := time.Now()
	decimalValue := decimal.NewFromFloat(10.00)

	o := &orderentity.Order{
		Entity: entity.Entity{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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
				TotalPayable:  decimal.NewFromFloat(150.45),
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
						StartAt:    &now,
						PendingAt:  &now,
						StartedAt:  &now,
						ReadyAt:    &now,
						CanceledAt: &now,
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

	out, err := FormatOrder(o)
	assert.NoError(t, err)
	if err := os.WriteFile("printer_output.txt", out, 0644); err != nil {
		t.Fatalf("failed to write printer buffer to file: %v", err)
	}
	s := string(out)
	assert.True(t, strings.HasPrefix(s, escInit))
	assert.True(t, strings.Contains(s, fmt.Sprintf("PEDIDO %d%s", o.OrderNumber, newline)))
	assert.True(t, strings.HasSuffix(s, escCut))
	// Delivery time logs
	assert.True(t, strings.Contains(s, fmt.Sprintf("Pendente: %s%s", now.Format("02/01/2006 15:04"), newline)))
	assert.True(t, strings.Contains(s, fmt.Sprintf("Despachado: %s%s", now.Format("02/01/2006 15:04"), newline)))
	// Delivery info
	assert.True(t, strings.Contains(s, fmt.Sprintf("Cliente: %s%s", o.Delivery.Client.Name, newline)))
	assert.True(t, strings.Contains(s, fmt.Sprintf("Endereço: %s, %s%s", o.Delivery.Address.Street, o.Delivery.Address.Number, newline)))
	assert.True(t, strings.Contains(s, fmt.Sprintf("Complemento: %s%s", o.Delivery.Address.Complement, newline)))
	assert.True(t, strings.Contains(s, fmt.Sprintf("Ref: %s%s", o.Delivery.Address.Reference, newline)))
	assert.True(t, strings.Contains(s, fmt.Sprintf("Bairro: %s%s", o.Delivery.Address.Neighborhood, newline)))
	assert.True(t, strings.Contains(s, fmt.Sprintf("Cidade: %s - %s%s", o.Delivery.Address.City, o.Delivery.Address.UF, newline)))
	assert.True(t, strings.Contains(s, fmt.Sprintf("CEP: %s%s", o.Delivery.Address.Cep, newline)))
	assert.True(t, strings.Contains(s, fmt.Sprintf("Motoboy: %s%s", o.Delivery.Driver.Employee.User.Name, newline)))
	// Delivery charges
	assert.True(t, strings.Contains(s, fmt.Sprintf("Taxa entrega: %7.2f%s", decimalValue.InexactFloat64(), newline)))
	assert.True(t, strings.Contains(s, fmt.Sprintf("Troco entrega: %7.2f%s", decimalValue.InexactFloat64(), newline)))
	assert.True(t, strings.Contains(s, fmt.Sprintf("Pagamento entrega: %s%s", o.Delivery.PaymentMethod, newline)))
	// Order details
	assert.True(t, strings.Contains(s, fmt.Sprintf("Observação: %s%s", o.Observation, newline)))
	assert.True(t, strings.Contains(s, fmt.Sprintf("Itens:%36.0f%s", o.QuantityItems, newline)))
	assert.True(t, strings.Contains(s, fmt.Sprintf("Pago:%31.2f%s", decimal.NewFromFloat(120.45).InexactFloat64(), newline)))
	assert.True(t, strings.Contains(s, fmt.Sprintf("Troco:%31.2f%s", decimal.NewFromFloat(30.45).InexactFloat64(), newline)))
	// Payments
	assert.True(t, strings.Contains(s, fmt.Sprintf("%s:%7.2f%s", o.Payments[0].Method, decimal.NewFromFloat(30.45).InexactFloat64(), newline)))
	assert.True(t, strings.Contains(s, fmt.Sprintf("%s:%7.2f%s", o.Payments[1].Method, decimal.NewFromFloat(20).InexactFloat64(), newline)))
	// Total payable
	assert.True(t, strings.Contains(s, fmt.Sprintf("TOTAL:%31.2f%s", decimal.NewFromFloat(150.45).InexactFloat64(), newline)))
}
