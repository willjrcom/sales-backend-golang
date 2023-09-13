package orderentity

type StatusOrder string

const (
	OrderStatusStaging  StatusOrder = "Staging"
	OrderStatusPending  StatusOrder = "Pending"
	OrderStatusReady    StatusOrder = "Ready"
	OrderStatusShipped  StatusOrder = "Shipped"
	OrderStatusFinished StatusOrder = "Finished"
	OrderStatusCanceled StatusOrder = "Canceled"
	OrderStatusArchived StatusOrder = "Archived"
)

func GetAllStatus() []StatusOrder {
	return []StatusOrder{
		OrderStatusStaging,
		OrderStatusPending,
		OrderStatusReady,
		OrderStatusShipped,
		OrderStatusFinished,
		OrderStatusCanceled,
		OrderStatusArchived,
	}
}
