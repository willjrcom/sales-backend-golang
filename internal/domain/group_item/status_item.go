package groupitementity

type StatusItem string

const (
	StatusGroupStaging  StatusItem = "Staging"
	StatusGroupPending  StatusItem = "Pending"
	StatusGroupReady    StatusItem = "Ready"
	StatusGroupCanceled StatusItem = "Canceled"
	StatusGroupFinished StatusItem = "Finished"
)

func GetAllStatus() []StatusItem {
	return []StatusItem{
		StatusGroupStaging,
		StatusGroupPending,
		StatusGroupReady,
		StatusGroupCanceled,
		StatusGroupFinished,
	}
}
