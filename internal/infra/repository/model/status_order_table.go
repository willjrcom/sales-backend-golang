package model

type StatusOrderTable string

const (
	OrderTableStatusStaging StatusOrderTable = "Staging"
	OrderTableStatusPending StatusOrderTable = "Pending"
	OrderTableStatusClosed  StatusOrderTable = "Closed"
)
