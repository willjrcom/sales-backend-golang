package orderentity

type StatusPickupOrder string

const (
	PickupOrderStatusStaging StatusPickupOrder = "Staging"
	PickupOrderStatusPending StatusPickupOrder = "Pending"
	PickupOrderStatusReady   StatusPickupOrder = "Ready"
)

func GetAllPickupStatus() []StatusPickupOrder {
	return []StatusPickupOrder{
		PickupOrderStatusStaging,
		PickupOrderStatusPending,
		PickupOrderStatusReady,
	}
}
