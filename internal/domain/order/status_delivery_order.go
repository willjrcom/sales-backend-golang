package orderentity

type StatusDeliveryOrder string

const (
	DeliveryOrderStatusStaging   StatusDeliveryOrder = "Staging"
	DeliveryOrderStatusPending   StatusDeliveryOrder = "Pending"
	DeliveryOrderStatusReady     StatusDeliveryOrder = "Ready"
	DeliveryOrderStatusShipped   StatusDeliveryOrder = "Shipped"
	DeliveryOrderStatusDelivered StatusDeliveryOrder = "Delivered"
)

func GetAllDeliveryStatus() []StatusDeliveryOrder {
	return []StatusDeliveryOrder{
		DeliveryOrderStatusStaging,
		DeliveryOrderStatusPending,
		DeliveryOrderStatusReady,
		DeliveryOrderStatusShipped,
		DeliveryOrderStatusDelivered,
	}
}
