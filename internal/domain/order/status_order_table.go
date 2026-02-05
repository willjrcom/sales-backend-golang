package orderentity

type StatusOrderTable string

const (
	OrderTableStatusStaging   StatusOrderTable = "Staging"
	OrderTableStatusPending   StatusOrderTable = "Pending"
	OrderTableStatusClosed    StatusOrderTable = "Closed"
	OrderTableStatusCancelled StatusOrderTable = "Cancelled"
)

func GetAllTableStatus() []StatusOrderTable {
	return []StatusOrderTable{
		OrderTableStatusStaging,
		OrderTableStatusPending,
		OrderTableStatusClosed,
		OrderTableStatusCancelled,
	}
}
