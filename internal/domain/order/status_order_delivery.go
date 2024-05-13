package orderentity

type StatusOrderDelivery string

const (
	OrderDeliveryStatusStaging   StatusOrderDelivery = "Staging"
	OrderDeliveryStatusPending   StatusOrderDelivery = "Pending"
	OrderDeliveryStatusReady     StatusOrderDelivery = "Ready"
	OrderDeliveryStatusShipped   StatusOrderDelivery = "Shipped"
	OrderDeliveryStatusDelivered StatusOrderDelivery = "Delivered"
)

func GetAllDeliveryStatus() []StatusOrderDelivery {
	return []StatusOrderDelivery{
		OrderDeliveryStatusStaging,
		OrderDeliveryStatusPending,
		OrderDeliveryStatusReady,
		OrderDeliveryStatusShipped,
		OrderDeliveryStatusDelivered,
	}
}
