package itementity

type StatusItem string

const (
	StatusItemPending  StatusItem = "Pending"
	StatusItemReady    StatusItem = "Ready"
	StatusItemCanceled StatusItem = "Canceled"
	StatusItemFinished StatusItem = "Finished"
)

func GetAllStatus() []StatusItem {
	return []StatusItem{
		StatusItemPending,
		StatusItemReady,
		StatusItemCanceled,
		StatusItemFinished,
	}
}
