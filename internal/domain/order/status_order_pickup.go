package orderentity

type StatusOrderPickup string

const (
	OrderPickupStatusStaging  StatusOrderPickup = "Staging"
	OrderPickupStatusPending  StatusOrderPickup = "Pending"
	OrderPickupStatusReady    StatusOrderPickup = "Ready"
	OrderPickupStatusCanceled StatusOrderPickup = "Canceled"
)

func GetAllPickupStatus() []StatusOrderPickup {
	return []StatusOrderPickup{
		OrderPickupStatusStaging,
		OrderPickupStatusPending,
		OrderPickupStatusReady,
		OrderPickupStatusCanceled,
	}
}
