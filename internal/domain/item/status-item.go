package itementity

type StatusItem string

const (
	StatusItemPending  StatusItem = "Pending"
	StatusItemReady    StatusItem = "Ready"
	StatusItemCanceled StatusItem = "Ready"
)

func getAllStatus() []StatusItem {
	return []StatusItem{
		StatusItemPending,
		StatusItemReady,
		StatusItemCanceled,
	}
}
