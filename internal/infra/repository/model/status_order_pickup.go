package model

type StatusOrderPickup string

const (
	OrderPickupStatusStaging StatusOrderPickup = "Staging"
	OrderPickupStatusPending StatusOrderPickup = "Pending"
	OrderPickupStatusReady   StatusOrderPickup = "Ready"
)
