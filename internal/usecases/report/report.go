package reportusecases

import (
	"context"
	"time"

	reportdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/report"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/report"
)

// Service wraps internal report generation logic.
type Service struct {
	reportSvc *report.ReportService
}

// TopTables returns the top 10 most used tables in the given period.
func (s *Service) TopTables(ctx context.Context, req *reportdto.TopTablesRequest) ([]reportdto.TopTablesResponse, error) {
	data, err := s.reportSvc.TopTables(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.TopTablesResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.TopTablesResponse{TableID: d.TableID, Count: d.Count}
	}
	return resp, nil
}

// NewService creates a new report usecase service.
func NewService(reportSvc *report.ReportService) *Service {
	return &Service{reportSvc: reportSvc}
}

// SalesTotalByDay returns total sales per day in the period.
func (s *Service) SalesTotalByDay(ctx context.Context, req *reportdto.SalesTotalByDayRequest) ([]reportdto.SalesByDayResponse, error) {
	data, err := s.reportSvc.SalesTotalByDay(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.SalesByDayResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.SalesByDayResponse{Day: d.Day, Total: d.Total}
	}
	return resp, nil
}

// RevenueCumulativeByMonth returns cumulative revenue by month.
func (s *Service) RevenueCumulativeByMonth(ctx context.Context, req *reportdto.RevenueCumulativeByMonthRequest) ([]reportdto.RevenueCumulativeByMonthResponse, error) {
	data, err := s.reportSvc.RevenueCumulativeByMonth(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.RevenueCumulativeByMonthResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.RevenueCumulativeByMonthResponse{Month: d.Month, Cumulative: d.Cumulative}
	}
	return resp, nil
}

// SalesByHour returns total sales per hour for a specific day.
func (s *Service) SalesByHour(ctx context.Context, req *reportdto.SalesByHourRequest) ([]reportdto.SalesByHourResponse, error) {
	data, err := s.reportSvc.SalesByHour(ctx, req.Day)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.SalesByHourResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.SalesByHourResponse{Hour: d.Hour, Total: d.Total}
	}
	return resp, nil
}

// SalesByChannel returns total sales per channel.
func (s *Service) SalesByChannel(ctx context.Context, req *reportdto.SalesByChannelRequest) ([]reportdto.SalesByChannelResponse, error) {
	data, err := s.reportSvc.SalesByChannel(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.SalesByChannelResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.SalesByChannelResponse{Channel: d.Channel, Total: d.Total}
	}
	return resp, nil
}

// AvgTicketByDay returns average ticket per day.
func (s *Service) AvgTicketByDay(ctx context.Context, req *reportdto.AvgTicketByDayRequest) ([]reportdto.AvgTicketByDayResponse, error) {
	data, err := s.reportSvc.AvgTicketByDay(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.AvgTicketByDayResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.AvgTicketByDayResponse{Day: d.Key.(time.Time), Avg: d.Avg}
	}
	return resp, nil
}

// AvgTicketByChannel returns average ticket per channel.
func (s *Service) AvgTicketByChannel(ctx context.Context, req *reportdto.AvgTicketByChannelRequest) ([]reportdto.AvgTicketByChannelResponse, error) {
	data, err := s.reportSvc.AvgTicketByChannel(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.AvgTicketByChannelResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.AvgTicketByChannelResponse{Channel: d.Key.(string), Avg: d.Avg}
	}
	return resp, nil
}

// ProductsSoldByDay returns total quantity of items sold per day.
func (s *Service) ProductsSoldByDay(ctx context.Context, req *reportdto.ProductsSoldByDayRequest) ([]reportdto.ProductsSoldByDayResponse, error) {
	data, err := s.reportSvc.ProductsSoldByDay(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.ProductsSoldByDayResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.ProductsSoldByDayResponse{Day: d.Day, Quantity: d.Quantity}
	}
	return resp, nil
}

// TopProducts returns top N products by quantity sold.
func (s *Service) TopProducts(ctx context.Context, req *reportdto.TopProductsRequest) ([]reportdto.TopProductsResponse, error) {
	data, err := s.reportSvc.TopProducts(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.TopProductsResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.TopProductsResponse{Name: d.Name, Quantity: d.Quantity}
	}
	return resp, nil
}

// SalesByCategory returns sales grouped by product category.
func (s *Service) SalesByCategory(ctx context.Context, req *reportdto.SalesByCategoryRequest) ([]reportdto.SalesByCategoryResponse, error) {
	data, err := s.reportSvc.SalesByCategory(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.SalesByCategoryResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.SalesByCategoryResponse{Category: d.Category, Quantity: d.Quantity}
	}
	return resp, nil
}

// ClientsRegisteredByDay returns count of clients registered per day.
func (s *Service) ClientsRegisteredByDay(ctx context.Context, req *reportdto.ClientsRegisteredByDayRequest) ([]reportdto.ClientsRegisteredByDayResponse, error) {
	data, err := s.reportSvc.ClientsRegisteredByDay(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.ClientsRegisteredByDayResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.ClientsRegisteredByDayResponse{Day: d.Day, Count: d.Count}
	}
	return resp, nil
}

// NewVsRecurringClients returns count of new vs recurring clients.
func (s *Service) NewVsRecurringClients(ctx context.Context, req *reportdto.NewVsRecurringClientsRequest) ([]reportdto.NewVsRecurringClientsResponse, error) {
	data, err := s.reportSvc.NewVsRecurringClients(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.NewVsRecurringClientsResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.NewVsRecurringClientsResponse{Type: d.Type, Count: d.Count}
	}
	return resp, nil
}

// OrdersByStatus returns count of orders per status.
func (s *Service) OrdersByStatus(ctx context.Context, req *reportdto.OrdersByStatusRequest) ([]reportdto.OrdersByStatusResponse, error) {
	data, err := s.reportSvc.OrdersByStatus(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.OrdersByStatusResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.OrdersByStatusResponse{Status: d.Status, Count: d.Count}
	}
	return resp, nil
}

// AvgProcessStepDurationByRule returns average duration per process rule.
func (s *Service) AvgProcessStepDurationByRule(ctx context.Context, req *reportdto.AvgProcessStepDurationRequest) ([]reportdto.AvgProcessStepDurationResponse, error) {
	data, err := s.reportSvc.AvgProcessStepDurationByRule(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.AvgProcessStepDurationResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.AvgProcessStepDurationResponse{ProcessRuleName: d.ProcessRuleName, AvgSeconds: d.AvgSeconds}
	}
	return resp, nil
}

// CancellationRate returns the cancellation rate of orders.
func (s *Service) CancellationRate(ctx context.Context, req *reportdto.CancellationRateRequest) (reportdto.CancellationRateResponse, error) {
	d, err := s.reportSvc.CancellationRate(ctx)
	if err != nil {
		return reportdto.CancellationRateResponse{}, err
	}
	return reportdto.CancellationRateResponse{Rate: d.Rate}, nil
}

// CurrentQueueLength returns the number of items currently in queue.
func (s *Service) CurrentQueueLength(ctx context.Context, req *reportdto.CurrentQueueLengthRequest) (reportdto.CurrentQueueLengthResponse, error) {
	d, err := s.reportSvc.CurrentQueueLength(ctx)
	if err != nil {
		return reportdto.CurrentQueueLengthResponse{}, err
	}
	return reportdto.CurrentQueueLengthResponse{Length: d.Length}, nil
}

// AvgDeliveryTimeByDriver returns average delivery time per driver.
func (s *Service) AvgDeliveryTimeByDriver(ctx context.Context, req *reportdto.AvgDeliveryTimeByDriverRequest) ([]reportdto.AvgDeliveryTimeByDriverResponse, error) {
	data, err := s.reportSvc.AvgDeliveryTimeByDriver(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.AvgDeliveryTimeByDriverResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.AvgDeliveryTimeByDriverResponse{DriverID: d.DriverID, AvgSeconds: d.AvgSeconds}
	}
	return resp, nil
}

// DeliveriesPerDriver returns number of deliveries per driver.
func (s *Service) DeliveriesPerDriver(ctx context.Context, req *reportdto.DeliveriesPerDriverRequest) ([]reportdto.DeliveriesPerDriverResponse, error) {
	data, err := s.reportSvc.DeliveriesPerDriver(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.DeliveriesPerDriverResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.DeliveriesPerDriverResponse{DriverID: d.DriverID, Count: d.Count}
	}
	return resp, nil
}

// OrdersPerTable returns number of orders per table.
func (s *Service) OrdersPerTable(ctx context.Context, req *reportdto.OrdersPerTableRequest) ([]reportdto.OrdersPerTableResponse, error) {
	data, err := s.reportSvc.OrdersPerTable(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.OrdersPerTableResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.OrdersPerTableResponse{TableID: d.TableID, Count: d.Count}
	}
	return resp, nil
}

// SalesByShift returns total sales per shift.
func (s *Service) SalesByShift(ctx context.Context, req *reportdto.SalesByShiftRequest) ([]reportdto.SalesByShiftResponse, error) {
	data, err := s.reportSvc.SalesByShift(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.SalesByShiftResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.SalesByShiftResponse{ShiftID: d.ShiftID, Total: d.Total}
	}
	return resp, nil
}

// PaymentsByMethod returns total payments by payment method.
func (s *Service) PaymentsByMethod(ctx context.Context, req *reportdto.PaymentsByMethodRequest) ([]reportdto.PaymentsByMethodResponse, error) {
	data, err := s.reportSvc.PaymentsByMethod(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.PaymentsByMethodResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.PaymentsByMethodResponse{Method: d.Method, Total: d.Total}
	}
	return resp, nil
}

// EmployeePaymentsReport returns sum of employee payments by employee.
func (s *Service) EmployeePaymentsReport(ctx context.Context, req *reportdto.EmployeePaymentsReportRequest) ([]reportdto.EmployeePaymentsReportResponse, error) {
	data, err := s.reportSvc.EmployeePaymentsReport(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.EmployeePaymentsReportResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.EmployeePaymentsReportResponse{EmployeeID: d.EmployeeID, Total: d.Total}
	}
	return resp, nil
}

// SalesByPlace returns total sales per place.
func (s *Service) SalesByPlace(ctx context.Context, req *reportdto.SalesByPlaceRequest) ([]reportdto.SalesByPlaceResponse, error) {
	data, err := s.reportSvc.SalesByPlace(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.SalesByPlaceResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.SalesByPlaceResponse{Place: d.Place, Total: d.Total}
	}
	return resp, nil
}

// SalesBySize returns total quantity sold per size.
func (s *Service) SalesBySize(ctx context.Context, req *reportdto.SalesBySizeRequest) ([]reportdto.SalesBySizeResponse, error) {
	data, err := s.reportSvc.SalesBySize(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.SalesBySizeResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.SalesBySizeResponse{Size: d.Size, Quantity: d.Quantity}
	}
	return resp, nil
}

// AdditionalItemsSold returns total quantity of additional items.
func (s *Service) AdditionalItemsSold(ctx context.Context, req *reportdto.AdditionalItemsRequest) ([]reportdto.AdditionalItemsResponse, error) {
	data, err := s.reportSvc.AdditionalItemsSold(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.AdditionalItemsResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.AdditionalItemsResponse{Name: d.Name, Quantity: d.Quantity}
	}
	return resp, nil
}

// AvgPickupTime returns average pickup wait time.
func (s *Service) AvgPickupTime(ctx context.Context, req *reportdto.AvgPickupTimeRequest) (reportdto.AvgPickupTimeResponse, error) {
	d, err := s.reportSvc.AvgPickupTime(ctx, req.Start, req.End)
	if err != nil {
		return reportdto.AvgPickupTimeResponse{}, err
	}
	return reportdto.AvgPickupTimeResponse{AvgSeconds: d.AvgSeconds}, nil
}

// GroupItemsByStatus returns count of group items by status.
func (s *Service) GroupItemsByStatus(ctx context.Context, req *reportdto.GroupItemsByStatusRequest) ([]reportdto.GroupItemsByStatusResponse, error) {
	data, err := s.reportSvc.GroupItemsByStatus(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.GroupItemsByStatusResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.GroupItemsByStatusResponse{Status: d.Status, Count: d.Count}
	}
	return resp, nil
}

// DeliveriesByCep returns number of deliveries per CEP.
func (s *Service) DeliveriesByCep(ctx context.Context, req *reportdto.DeliveriesByCepRequest) ([]reportdto.DeliveriesByCepResponse, error) {
	data, err := s.reportSvc.DeliveriesByCep(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.DeliveriesByCepResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.DeliveriesByCepResponse{Cep: d.Cep, Count: d.Count}
	}
	return resp, nil
}

// ProcessedCountByRule returns processed counts by rule.
func (s *Service) ProcessedCountByRule(ctx context.Context, req *reportdto.ProcessedCountByRuleRequest) ([]reportdto.ProcessedCountByRuleResponse, error) {
	data, err := s.reportSvc.ProcessedCountByRule(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.ProcessedCountByRuleResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.ProcessedCountByRuleResponse{RuleID: d.RuleID, Count: d.Count}
	}
	return resp, nil
}

// DailySales returns summary metrics for a specific day.
func (s *Service) DailySales(ctx context.Context, req *reportdto.DailySalesRequest) (reportdto.DailySalesResponse, error) {
	data, err := s.reportSvc.DailySales(ctx, req.Day)
	if err != nil {
		return reportdto.DailySalesResponse{}, err
	}
	return reportdto.DailySalesResponse{TotalOrders: data.TotalOrders, TotalSales: data.TotalSales}, nil
}

// AvgQueueDuration returns average duration of all queues.
func (s *Service) AvgQueueDuration(ctx context.Context, req *reportdto.AvgQueueDurationRequest) (reportdto.AvgQueueDurationResponse, error) {
	d, err := s.reportSvc.AvgQueueDuration(ctx)
	if err != nil {
		return reportdto.AvgQueueDurationResponse{}, err
	}
	return reportdto.AvgQueueDurationResponse{AvgSeconds: d.AvgSeconds}, nil
}

// AvgProcessDurationByProduct returns average process duration by product.
func (s *Service) AvgProcessDurationByProduct(ctx context.Context, req *reportdto.AvgProcessDurationByProductRequest) ([]reportdto.AvgProcessDurationByProductResponse, error) {
	data, err := s.reportSvc.AvgProcessDurationByProduct(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.AvgProcessDurationByProductResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.AvgProcessDurationByProductResponse{
			ProductID:   d.ProductID,
			ProductName: d.ProductName,
			AvgSeconds:  d.AvgSeconds,
		}
	}
	return resp, nil
}

// TotalQueueTimeByGroupItem returns total sum of queue durations per group item.
func (s *Service) TotalQueueTimeByGroupItem(ctx context.Context, req *reportdto.TotalQueueTimeByGroupItemRequest) ([]reportdto.TotalQueueTimeByGroupItemResponse, error) {
	data, err := s.reportSvc.TotalQueueTimeByGroupItem(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.TotalQueueTimeByGroupItemResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.TotalQueueTimeByGroupItemResponse{
			GroupItemID:  d.Name,
			TotalSeconds: d.TotalSeconds,
		}
	}
	return resp, nil
}

// ProductProfitability returns profitability analysis per product.
func (s *Service) ProductProfitability(ctx context.Context, req *reportdto.ProductProfitabilityRequest) ([]reportdto.ProductProfitabilityResponse, error) {
	data, err := s.reportSvc.ProductProfitability(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.ProductProfitabilityResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.ProductProfitabilityResponse{
			ProductID:    d.ProductID,
			ProductName:  d.ProductName,
			TotalSold:    d.TotalSold,
			TotalCost:    d.TotalCost,
			TotalRevenue: d.TotalRevenue,
			Profit:       d.Profit,
			ProfitMargin: d.ProfitMargin,
		}
	}
	return resp, nil
}

// CategoryProfitability returns profitability analysis per category.
func (s *Service) CategoryProfitability(ctx context.Context, req *reportdto.CategoryProfitabilityRequest) ([]reportdto.CategoryProfitabilityResponse, error) {
	data, err := s.reportSvc.CategoryProfitability(ctx, req.Start, req.End)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.CategoryProfitabilityResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.CategoryProfitabilityResponse{
			CategoryName: d.CategoryName,
			TotalSold:    d.TotalSold,
			TotalCost:    d.TotalCost,
			TotalRevenue: d.TotalRevenue,
			Profit:       d.Profit,
			ProfitMargin: d.ProfitMargin,
		}
	}
	return resp, nil
}

// LowProfitProducts returns products with low profit margin.
func (s *Service) LowProfitProducts(ctx context.Context, req *reportdto.LowProfitProductsRequest) ([]reportdto.LowProfitProductsResponse, error) {
	data, err := s.reportSvc.LowProfitProducts(ctx, req.MinMargin)
	if err != nil {
		return nil, err
	}
	resp := make([]reportdto.LowProfitProductsResponse, len(data))
	for i, d := range data {
		resp[i] = reportdto.LowProfitProductsResponse{
			ProductID:    d.ProductID,
			ProductName:  d.ProductName,
			Price:        d.Price,
			Cost:         d.Cost,
			ProfitMargin: d.ProfitMargin,
		}
	}
	return resp, nil
}

func (s *Service) DebugProducts(ctx context.Context) (interface{}, error) {
	return s.reportSvc.DebugProducts(ctx)
}

// OverallProfitability returns overall profitability metrics.
func (s *Service) OverallProfitability(ctx context.Context, req *reportdto.OverallProfitabilityRequest) (reportdto.OverallProfitabilityResponse, error) {
	data, err := s.reportSvc.OverallProfitability(ctx, req.Start, req.End)
	if err != nil {
		return reportdto.OverallProfitabilityResponse{}, err
	}
	return reportdto.OverallProfitabilityResponse{
		TotalRevenue: data.TotalRevenue,
		TotalCost:    data.TotalCost,
		TotalProfit:  data.TotalProfit,
		ProfitMargin: data.ProfitMargin,
	}, nil
}
