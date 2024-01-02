package orderdto

import orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"

type UpdatePaymentMethod struct {
	orderentity.PaymentOrder
}

func (u *UpdatePaymentMethod) UpdateModel(order *orderentity.Order) {
	order.Payment = &orderentity.PaymentOrder{}

	if u.TotalPaid != nil {
		order.Payment.TotalPaid = u.TotalPaid
	}
	if u.Change != nil {
		order.Payment.Change = u.Change
	}
	if u.Method != nil {
		order.Payment.Method = u.Method
	}
}
