package itementity

type StatusItem string

const (
	StatusItemPending  StatusItem = "Pending"
	StatusItemReady    StatusItem = "Ready"
	StatusItemCanceled StatusItem = "Ready"
)

func GetAllStatus() []StatusItem {
	return []StatusItem{
		StatusItemPending,
		StatusItemReady,
		StatusItemCanceled,
	}
}
