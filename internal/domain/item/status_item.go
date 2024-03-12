package itementity

type StatusItem string

const (
	StatusItemStaging  StatusItem = "Staging"
	StatusItemPending  StatusItem = "Pending"
	StatusItemStarted  StatusItem = "Started"
	StatusItemReady    StatusItem = "Ready"
	StatusItemCanceled StatusItem = "Canceled"
)

func GetAllStatus() []StatusItem {
	return []StatusItem{
		StatusItemStaging,
		StatusItemPending,
		StatusItemReady,
		StatusItemCanceled,
		StatusItemStarted,
	}
}
