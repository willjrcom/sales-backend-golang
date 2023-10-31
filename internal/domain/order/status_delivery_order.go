package orderentity

type StatusDeliveryOrder string

const (
	DeliveryOrderStatusPending   StatusDeliveryOrder = "Pending"
	DeliveryOrderStatusShipped   StatusDeliveryOrder = "Shipped"
	DeliveryOrderStatusDelivered StatusDeliveryOrder = "Delivered"
)

func GetAllDeliveryStatus() []StatusDeliveryOrder {
	return []StatusDeliveryOrder{
		DeliveryOrderStatusPending,
		DeliveryOrderStatusShipped,
		DeliveryOrderStatusDelivered,
	}
}
