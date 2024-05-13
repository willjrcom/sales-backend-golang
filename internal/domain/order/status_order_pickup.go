package orderentity

type StatusOrderPickup string

const (
	OrderPickupStatusStaging StatusOrderPickup = "Staging"
	OrderPickupStatusPending StatusOrderPickup = "Pending"
	OrderPickupStatusReady   StatusOrderPickup = "Ready"
)

func GetAllPickupStatus() []StatusOrderPickup {
	return []StatusOrderPickup{
		OrderPickupStatusStaging,
		OrderPickupStatusPending,
		OrderPickupStatusReady,
	}
}
