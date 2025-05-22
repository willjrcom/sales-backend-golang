package reportdto

import (
    "time"

    "github.com/shopspring/decimal"
)

// SalesTotalByDayRequest filters for sales total by day.
type SalesTotalByDayRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// SalesByDayResponse holds total sales for a day.
type SalesByDayResponse struct {
    Day   time.Time       `json:"day"`
    Total decimal.Decimal `json:"total"`
}

// RevenueCumulativeByMonthRequest filters for cumulative monthly revenue.
type RevenueCumulativeByMonthRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// RevenueCumulativeByMonthResponse holds cumulative revenue for a month.
type RevenueCumulativeByMonthResponse struct {
    Month      time.Time       `json:"month"`
    Cumulative decimal.Decimal `json:"cumulative"`
}

// SalesByHourRequest filters for hourly sales.
type SalesByHourRequest struct {
    Schema string    `json:"schema"`
    Day    time.Time `json:"day"`
}

// SalesByHourResponse holds sales total for an hour.
type SalesByHourResponse struct {
    Hour  int             `json:"hour"`
    Total decimal.Decimal `json:"total"`
}

// SalesByChannelRequest filters for sales by channel.
type SalesByChannelRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// SalesByChannelResponse holds sales total per channel.
type SalesByChannelResponse struct {
    Channel string          `json:"channel"`
    Total   decimal.Decimal `json:"total"`
}

// AvgTicketByDayRequest filters for average ticket by day.
type AvgTicketByDayRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// AvgTicketByDayResponse holds average ticket for a day.
type AvgTicketByDayResponse struct {
    Day time.Time       `json:"day"`
    Avg decimal.Decimal `json:"avg"`
}

// AvgTicketByChannelRequest filters for average ticket by channel.
type AvgTicketByChannelRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// AvgTicketByChannelResponse holds average ticket for a channel.
type AvgTicketByChannelResponse struct {
    Channel string          `json:"channel"`
    Avg     decimal.Decimal `json:"avg"`
}

// ProductsSoldByDayRequest filters for products sold by day.
type ProductsSoldByDayRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// ProductsSoldByDayResponse holds total products sold for a day.
type ProductsSoldByDayResponse struct {
    Day      time.Time `json:"day"`
    Quantity float64   `json:"quantity"`
}

// TopProductsRequest filters for top products.
type TopProductsRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// TopProductsResponse holds product name and quantity sold.
type TopProductsResponse struct {
    Name     string  `json:"name"`
    Quantity float64 `json:"quantity"`
}

// SalesByCategoryRequest filters for sales by category.
type SalesByCategoryRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// SalesByCategoryResponse holds category and quantity sold.
type SalesByCategoryResponse struct {
    Category string  `json:"category"`
    Quantity float64 `json:"quantity"`
}

// CurrentStockRequest filters for current stock by category.
type CurrentStockRequest struct {
    Schema string `json:"schema"`
}

// CurrentStockResponse holds category and current stock.
type CurrentStockResponse struct {
    Category string  `json:"category"`
    Quantity float64 `json:"quantity"`
}

// ClientsRegisteredByDayRequest filters for clients registered.
type ClientsRegisteredByDayRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// ClientsRegisteredByDayResponse holds date and count of new clients.
type ClientsRegisteredByDayResponse struct {
    Day   time.Time `json:"day"`
    Count int       `json:"count"`
}

// NewVsRecurringClientsRequest filters for new vs recurring clients.
type NewVsRecurringClientsRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// NewVsRecurringClientsResponse holds type and count.
type NewVsRecurringClientsResponse struct {
    Type  string `json:"type"`
    Count int    `json:"count"`
}

// OrdersByStatusRequest filters for orders by status.
type OrdersByStatusRequest struct {
    Schema string `json:"schema"`
}

// OrdersByStatusResponse holds status and count.
type OrdersByStatusResponse struct {
    Status string `json:"status"`
    Count  int    `json:"count"`
}

// AvgProcessStepDurationRequest filters for average process step duration.
type AvgProcessStepDurationRequest struct {
    Schema string `json:"schema"`
}

// AvgProcessStepDurationResponse holds process rule ID and average seconds.
type AvgProcessStepDurationResponse struct {
    ProcessRuleID string  `json:"process_rule_id"`
    AvgSeconds    float64 `json:"avg_seconds"`
}

// CancellationRateRequest filters for cancellation rate.
type CancellationRateRequest struct {
    Schema string `json:"schema"`
}

// CancellationRateResponse holds cancellation rate.
type CancellationRateResponse struct {
    Rate float64 `json:"rate"`
}

// CurrentQueueLengthRequest filters for current queue length.
type CurrentQueueLengthRequest struct {
    Schema string `json:"schema"`
}

// CurrentQueueLengthResponse holds length of queue.
type CurrentQueueLengthResponse struct {
    Length int `json:"length"`
}

// AvgDeliveryTimeByDriverRequest filters for avg delivery time.
type AvgDeliveryTimeByDriverRequest struct {
    Schema string `json:"schema"`
}

// AvgDeliveryTimeByDriverResponse holds driver ID and average seconds.
type AvgDeliveryTimeByDriverResponse struct {
    DriverID   string  `json:"driver_id"`
    AvgSeconds float64 `json:"avg_seconds"`
}

// DeliveriesPerDriverRequest filters for deliveries per driver.
type DeliveriesPerDriverRequest struct {
    Schema string `json:"schema"`
}

// DeliveriesPerDriverResponse holds driver ID and count.
type DeliveriesPerDriverResponse struct {
    DriverID string `json:"driver_id"`
    Count    int    `json:"count"`
}

// OrdersPerTableRequest filters for orders per table.
type OrdersPerTableRequest struct {
    Schema string `json:"schema"`
}

// OrdersPerTableResponse holds table ID and count.
type OrdersPerTableResponse struct {
    TableID string `json:"table_id"`
    Count   int    `json:"count"`
}

// SalesByShiftRequest filters for sales by shift.
type SalesByShiftRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// SalesByShiftResponse holds shift ID and total sales.
type SalesByShiftResponse struct {
    ShiftID string          `json:"shift_id"`
    Total   decimal.Decimal `json:"total"`
}

// PaymentsByMethodRequest filters for payments by method.
type PaymentsByMethodRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// PaymentsByMethodResponse holds method and total payments.
type PaymentsByMethodResponse struct {
    Method string          `json:"method"`
    Total  decimal.Decimal `json:"total"`
}

// EmployeePaymentsReportRequest filters for employee payments.
type EmployeePaymentsReportRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// EmployeePaymentsReportResponse holds employee ID and total paid.
type EmployeePaymentsReportResponse struct {
    EmployeeID string          `json:"employee_id"`
    Total      decimal.Decimal `json:"total"`
}
// SalesByPlaceRequest filters for sales by place.
type SalesByPlaceRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// SalesByPlaceResponse holds total sales per place.
type SalesByPlaceResponse struct {
    Place string          `json:"place"`
    Total decimal.Decimal `json:"total"`
}

// SalesBySizeRequest filters for sales by product size.
type SalesBySizeRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// SalesBySizeResponse holds quantity sold per size.
type SalesBySizeResponse struct {
    Size     string  `json:"size"`
    Quantity float64 `json:"quantity"`
}

// AdditionalItemsRequest filters for additional items sold.
type AdditionalItemsRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// AdditionalItemsResponse holds name and quantity of additionals.
type AdditionalItemsResponse struct {
    Name     string  `json:"name"`
    Quantity float64 `json:"quantity"`
}

// AvgPickupTimeRequest filters for average pickup wait time.
type AvgPickupTimeRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// AvgPickupTimeResponse holds average waiting time in seconds.
type AvgPickupTimeResponse struct {
    AvgSeconds float64 `json:"avg_seconds"`
}

// GroupItemsByStatusRequest filters for group item status counts.
type GroupItemsByStatusRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// GroupItemsByStatusResponse holds status and count.
type GroupItemsByStatusResponse struct {
    Status string `json:"status"`
    Count  int    `json:"count"`
}

// DeliveriesByCepRequest filters by delivery date range for heatmap.
type DeliveriesByCepRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// DeliveriesByCepResponse holds CEP and delivery count.
type DeliveriesByCepResponse struct {
    Cep   string `json:"cep"`
    Count int    `json:"count"`
}

// ProcessedCountByRuleRequest filters for processed count by rule.
type ProcessedCountByRuleRequest struct {
    Schema string    `json:"schema"`
    Start  time.Time `json:"start"`
    End    time.Time `json:"end"`
}

// ProcessedCountByRuleResponse holds rule ID and count.
type ProcessedCountByRuleResponse struct {
    RuleID string `json:"rule_id"`
    Count  int    `json:"count"`
}