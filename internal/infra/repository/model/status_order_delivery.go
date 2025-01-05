package model

type StatusOrderDelivery string

const (
	OrderDeliveryStatusStaging   StatusOrderDelivery = "Staging"
	OrderDeliveryStatusPending   StatusOrderDelivery = "Pending"
	OrderDeliveryStatusReady     StatusOrderDelivery = "Ready"
	OrderDeliveryStatusShipped   StatusOrderDelivery = "Shipped"
	OrderDeliveryStatusDelivered StatusOrderDelivery = "Delivered"
)
