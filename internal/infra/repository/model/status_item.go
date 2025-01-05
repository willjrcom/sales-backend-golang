package model

type StatusGroupItem string

const (
	StatusGroupStaging  StatusGroupItem = "Staging"
	StatusGroupPending  StatusGroupItem = "Pending"
	StatusGroupStarted  StatusGroupItem = "Started"
	StatusGroupReady    StatusGroupItem = "Ready"
	StatusGroupCanceled StatusGroupItem = "Canceled"
)
