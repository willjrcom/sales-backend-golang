package groupitementity

type StatusGroupItem string

const (
	StatusGroupStaging  StatusGroupItem = "Staging"
	StatusGroupPending  StatusGroupItem = "Pending"
	StatusGroupStarted  StatusGroupItem = "Started"
	StatusGroupReady    StatusGroupItem = "Ready"
	StatusGroupCanceled StatusGroupItem = "Canceled"
)

func GetAllStatus() []StatusGroupItem {
	return []StatusGroupItem{
		StatusGroupStaging,
		StatusGroupPending,
		StatusGroupStarted,
		StatusGroupReady,
		StatusGroupCanceled,
	}
}
