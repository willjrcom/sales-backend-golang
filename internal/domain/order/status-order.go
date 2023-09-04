package orderentity

type StatusOrder string

const (
	OrderStatusStaging  StatusOrder = "Staging"
	OrderStatusPending  StatusOrder = "Pending"
	OrderStatusReady    StatusOrder = "Ready"
	OrderStatusShipped  StatusOrder = "Shipped"
	OrderStatusFinished StatusOrder = "Finished"
	OrderStatusCanceled StatusOrder = "Canceled"
)

func getAllStatus() []StatusOrder {
	return []StatusOrder{
		OrderStatusStaging,
		OrderStatusPending,
		OrderStatusReady,
		OrderStatusShipped,
		OrderStatusFinished,
		OrderStatusCanceled,
	}
}
