package orderentity

type StatusOrderTable string

const (
	OrderTableStatusStaging StatusOrderTable = "Staging"
	OrderTableStatusPending StatusOrderTable = "Pending"
	OrderTableStatusClosed  StatusOrderTable = "Closed"
)

func GetAllTableStatus() []StatusOrderTable {
	return []StatusOrderTable{
		OrderTableStatusStaging,
		OrderTableStatusPending,
		OrderTableStatusClosed,
	}
}
