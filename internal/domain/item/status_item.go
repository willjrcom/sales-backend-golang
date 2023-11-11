package itementity

type StatusItem string

const (
	StatusItemStaging  StatusItem = "Staging"
	StatusItemPending  StatusItem = "Pending"
	StatusItemReady    StatusItem = "Ready"
	StatusItemCanceled StatusItem = "Canceled"
	StatusItemFinished StatusItem = "Finished"
)

func GetAllStatus() []StatusItem {
	return []StatusItem{
		StatusItemStaging,
		StatusItemPending,
		StatusItemReady,
		StatusItemCanceled,
		StatusItemFinished,
	}
}
