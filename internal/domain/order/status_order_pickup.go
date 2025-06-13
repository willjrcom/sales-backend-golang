package orderentity

type StatusOrderPickup string

const (
	OrderPickupStatusStaging   StatusOrderPickup = "Staging"
	OrderPickupStatusPending   StatusOrderPickup = "Pending"
	OrderPickupStatusReady     StatusOrderPickup = "Ready"
	OrderPickupStatusDelivered StatusOrderPickup = "Delivered"
	OrderPickupStatusCanceled  StatusOrderPickup = "Canceled"
)

func GetAllPickupStatus() []StatusOrderPickup {
	return []StatusOrderPickup{
		OrderPickupStatusStaging,
		OrderPickupStatusPending,
		OrderPickupStatusReady,
		OrderPickupStatusDelivered,
		OrderPickupStatusCanceled,
	}
}
