package orderentity

type StatusPickupOrder string

const (
	PickupOrderStatusPending  StatusPickupOrder = "Pending"
	PickupOrderStatusReady    StatusPickupOrder = "Ready"
	PickupOrderStatusPickedup StatusPickupOrder = "Picked up"
)

func GetAllPickupStatus() []StatusPickupOrder {
	return []StatusPickupOrder{
		PickupOrderStatusPending,
		PickupOrderStatusReady,
		PickupOrderStatusPickedup,
	}
}
