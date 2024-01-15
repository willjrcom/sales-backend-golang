package orderentity

type StatusOrder string

const (
	OrderStatusStaging  StatusOrder = "Staging"
	OrderStatusPending  StatusOrder = "Pending"
	OrderStatusFinished StatusOrder = "Finished"
	OrderStatusCanceled StatusOrder = "Canceled"
	OrderStatusArchived StatusOrder = "Archived"
)

func GetAllOrderStatus() []StatusOrder {
	return []StatusOrder{
		OrderStatusStaging,
		OrderStatusPending,
		OrderStatusFinished,
		OrderStatusCanceled,
		OrderStatusArchived,
	}
}
