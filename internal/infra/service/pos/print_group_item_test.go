package pos

import (
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

func Test_FormatGroupItem(t *testing.T) {
	now := time.Now().UTC()
	groupItem := orderentity.GroupItem{
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
	}

	out, err := FormatGroupItemKitchen(&groupItem, nil)
	assert.NoError(t, err)
	if err := os.WriteFile("printer_group_item.txt", out, 0644); err != nil {
		t.Fatalf("failed to write printer buffer to file: %v", err)
	}
}
