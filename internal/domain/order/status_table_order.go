package orderentity

type StatusTableOrder string

const (
	TableOrderStatusStaging StatusTableOrder = "Staging"
	TableOrderStatusPending StatusTableOrder = "Pending"
	TableOrderStatusClosed  StatusTableOrder = "Closed"
)

func GetAllTableStatus() []StatusTableOrder {
	return []StatusTableOrder{
		TableOrderStatusStaging,
		TableOrderStatusPending,
		TableOrderStatusClosed,
	}
}
