package orderentity

type StatusOrderTable string

const (
	OrderTableStatusStaging  StatusOrderTable = "Staging"
	OrderTableStatusPending  StatusOrderTable = "Pending"
	OrderTableStatusClosed   StatusOrderTable = "Closed"
	OrderTableStatusCanceled StatusOrderTable = "Canceled"
)

func GetAllTableStatus() []StatusOrderTable {
	return []StatusOrderTable{
		OrderTableStatusStaging,
		OrderTableStatusPending,
		OrderTableStatusClosed,
		OrderTableStatusCanceled,
	}
}
