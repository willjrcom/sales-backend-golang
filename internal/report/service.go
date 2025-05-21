package report

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
)

// ReportService provides methods to generate various sales and post-sales reports.
type ReportService struct {
	db *bun.DB
}

// NewReportService creates a new ReportService using the given DB connection.
func NewReportService(db *bun.DB) *ReportService {
	return &ReportService{db: db}
}

// SalesByDayDTO holds total payable grouped by day.
type SalesByDayDTO struct {
	Day   time.Time       `bun:"day"`
	Total decimal.Decimal `bun:"total"`
}

// SalesTotalByDay returns total sales (sum of total_payable) per day in the given period.
func (s *ReportService) SalesTotalByDay(ctx context.Context, schema string, start, end time.Time) ([]SalesByDayDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []SalesByDayDTO
	query := `
        SELECT date(created_at) AS day, SUM(total_payable) AS total
        FROM orders
        WHERE created_at BETWEEN ? AND ?
        GROUP BY day
        ORDER BY day`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CumulativeRevenueDTO holds cumulative revenue per month.
type CumulativeRevenueDTO struct {
	Month      time.Time       `bun:"mon"`
	Cumulative decimal.Decimal `bun:"cumulative_rev"`
}

// RevenueCumulativeByMonth returns cumulative monthly revenue.
func (s *ReportService) RevenueCumulativeByMonth(ctx context.Context, schema string, start, end time.Time) ([]CumulativeRevenueDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []CumulativeRevenueDTO
	query := `
        WITH M AS (
            SELECT date_trunc('month', created_at)::date AS mon, SUM(total_payable) AS rev
            FROM orders
            WHERE created_at BETWEEN ? AND ?
            GROUP BY mon
        )
        SELECT mon, SUM(rev) OVER (ORDER BY mon) AS cumulative_rev
        FROM M`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// HourlySalesDTO holds sales total per hour.
type HourlySalesDTO struct {
	Hour  int             `bun:"hr"`
	Total decimal.Decimal `bun:"total"`
}

// SalesByHour returns total sales per hour for a specific day.
func (s *ReportService) SalesByHour(ctx context.Context, schema string, day time.Time) ([]HourlySalesDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []HourlySalesDTO
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)
	query := `
        SELECT EXTRACT(hour FROM created_at)::int AS hr, SUM(total_payable) AS total
        FROM orders
        WHERE created_at >= ? AND created_at < ?
        GROUP BY hr
        ORDER BY hr`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ChannelSalesDTO holds sales per channel type.
type ChannelSalesDTO struct {
	Channel string          `bun:"channel"`
	Total   decimal.Decimal `bun:"total"`
}

// SalesByChannel returns sum of sales per channel: delivery, pickup, table.
func (s *ReportService) SalesByChannel(ctx context.Context, schema string, start, end time.Time) ([]ChannelSalesDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []ChannelSalesDTO
	query := `
        SELECT
            CASE
                WHEN d.id IS NOT NULL THEN 'delivery'
                WHEN t.id IS NOT NULL THEN 'table'
                WHEN p.id IS NOT NULL THEN 'pickup'
                ELSE 'unknown'
            END AS channel,
            SUM(o.total_payable) AS total
        FROM orders o
        LEFT JOIN order_deliveries d ON d.order_id = o.id
        LEFT JOIN order_tables t ON t.order_id = o.id
        LEFT JOIN order_pickups p ON p.order_id = o.id
        WHERE o.created_at BETWEEN ? AND ?
        GROUP BY channel`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AvgTicketDTO holds average ticket per grouping key.
type AvgTicketDTO struct {
	Key interface{}     `bun:"key"`
	Avg decimal.Decimal `bun:"avg_ticket"`
}

// AvgTicketByDay returns average order total per day.
func (s *ReportService) AvgTicketByDay(ctx context.Context, schema string, start, end time.Time) ([]AvgTicketDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []AvgTicketDTO
	query := `
        SELECT date(created_at) AS key, AVG(total_payable) AS avg_ticket
        FROM orders
        WHERE created_at BETWEEN ? AND ?
        GROUP BY key
        ORDER BY key`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AvgTicketByChannel returns average ticket per channel.
func (s *ReportService) AvgTicketByChannel(ctx context.Context, schema string, start, end time.Time) ([]AvgTicketDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []AvgTicketDTO
	query := `
        SELECT
            CASE
                WHEN d.id IS NOT NULL THEN 'delivery'
                WHEN t.id IS NOT NULL THEN 'table'
                WHEN p.id IS NOT NULL THEN 'pickup'
                ELSE 'unknown'
            END AS key,
            AVG(o.total_payable) AS avg_ticket
        FROM orders o
        LEFT JOIN order_deliveries d ON d.order_id = o.id
        LEFT JOIN order_tables t ON t.order_id = o.id
        LEFT JOIN order_pickups p ON p.order_id = o.id
        WHERE o.created_at BETWEEN ? AND ?
        GROUP BY key`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ProductsSoldByDayDTO holds daily sum of items sold.
type ProductsSoldByDayDTO struct {
	Day      time.Time `bun:"day"`
	Quantity float64   `bun:"quantity"`
}

// ProductsSoldByDay returns sum of quantity_items per day.
func (s *ReportService) ProductsSoldByDay(ctx context.Context, schema string, start, end time.Time) ([]ProductsSoldByDayDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []ProductsSoldByDayDTO
	query := `
        SELECT date(created_at) AS day, SUM(quantity_items) AS quantity
        FROM orders
        WHERE created_at BETWEEN ? AND ?
        GROUP BY day
        ORDER BY day`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// TopProductsDTO holds top products by quantity.
type TopProductsDTO struct {
	Name     string  `bun:"name"`
	Quantity float64 `bun:"quantity"`
}

// TopProducts returns top 10 products by sold quantity.
func (s *ReportService) TopProducts(ctx context.Context, schema string, start, end time.Time) ([]TopProductsDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []TopProductsDTO
	query := `
        SELECT p.name, SUM(i.quantity) AS quantity
        FROM order_items i
        JOIN order_group_items g ON i.group_item_id = g.id
        JOIN orders o ON g.order_id = o.id
        JOIN products p ON p.id = i.product_id
        WHERE o.created_at BETWEEN ? AND ?
        GROUP BY p.name
        ORDER BY quantity DESC
        LIMIT 10`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SalesByCategoryDTO holds sales by product category.
type SalesByCategoryDTO struct {
	Category string  `bun:"category"`
	Quantity float64 `bun:"quantity"`
}

// SalesByCategory returns sum of quantities by product category.
func (s *ReportService) SalesByCategory(ctx context.Context, schema string, start, end time.Time) ([]SalesByCategoryDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []SalesByCategoryDTO
	query := `
        SELECT pc.name AS category, SUM(i.quantity) AS quantity
        FROM order_items i
        JOIN product_categories pc ON pc.id = i.category_id
        WHERE i.created_at BETWEEN ? AND ?
        GROUP BY pc.name
        ORDER BY category`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CurrentStockDTO holds current stock by category.
type CurrentStockDTO struct {
	Category string  `bun:"category"`
	Quantity float64 `bun:"quantity"`
}

// CurrentStockByCategory returns current stock level per product category.
func (s *ReportService) CurrentStockByCategory(ctx context.Context, schema string) ([]CurrentStockDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []CurrentStockDTO
	query := `
        SELECT pc.name AS category, q.quantity
        FROM quantities q
        JOIN product_categories pc ON pc.id = q.category_id`
	if err := s.db.NewRaw(query).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: implement MovimentacaoEstoque (requires inventory log table)

// ClientsRegisteredDTO holds count of new clients by day.
type ClientsRegisteredDTO struct {
	Day   time.Time `bun:"day"`
	Count int       `bun:"count"`
}

// ClientsRegisteredByDay returns the number of clients registered per day.
func (s *ReportService) ClientsRegisteredByDay(ctx context.Context, schema string, start, end time.Time) ([]ClientsRegisteredDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []ClientsRegisteredDTO
	query := `
        SELECT date(created_at) AS day, COUNT(*) AS count
        FROM clients
        WHERE created_at BETWEEN ? AND ?
        GROUP BY day
        ORDER BY day`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// NewVsRecurringDTO holds ratio of new vs recurring clients.
type NewVsRecurringDTO struct {
	Type  string `bun:"type"`
	Count int    `bun:"count"`
}

// NewVsRecurringClients returns counts of new vs recurring clients for deliveries in period.
func (s *ReportService) NewVsRecurringClients(ctx context.Context, schema string, start, end time.Time) ([]NewVsRecurringDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []NewVsRecurringDTO
	query := `
        SELECT
            CASE WHEN c.created_at >= ? THEN 'new' ELSE 'recurring' END AS type,
            COUNT(DISTINCT d.client_id) AS count
        FROM order_deliveries d
        JOIN clients c ON c.id = d.client_id
        WHERE d.delivered_at BETWEEN ? AND ?
        GROUP BY type`
	if err := s.db.NewRaw(query, start, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: implement ContactConversionRate (requires order origin field)

// OrdersByStatusDTO holds count of orders by status.
type OrdersByStatusDTO struct {
	Status string `bun:"status"`
	Count  int    `bun:"count"`
}

// OrdersByStatus returns the number of orders per status.
func (s *ReportService) OrdersByStatus(ctx context.Context, schema string) ([]OrdersByStatusDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []OrdersByStatusDTO
	query := `
        SELECT status, COUNT(*) AS count
        FROM orders
        GROUP BY status`
	if err := s.db.NewRaw(query).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AvgProcessStepDTO holds average duration (in seconds) per process rule.
type AvgProcessStepDTO struct {
	ProcessRuleID string  `bun:"process_rule_id"`
	AvgSeconds    float64 `bun:"avg_seconds"`
}

// AvgProcessStepDurationByRule returns average time per process rule.
func (s *ReportService) AvgProcessStepDurationByRule(ctx context.Context, schema string) ([]AvgProcessStepDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []AvgProcessStepDTO
	query := `
        SELECT process_rule_id, AVG(EXTRACT(EPOCH FROM (finished_at - started_at))) AS avg_seconds
        FROM order_processes
        WHERE finished_at IS NOT NULL AND started_at IS NOT NULL
        GROUP BY process_rule_id`
	if err := s.db.NewRaw(query).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CancellationRateDTO holds cancellation rate metric.
type CancellationRateDTO struct {
	Rate float64 `bun:"rate"`
}

// CancellationRate returns percentage of orders canceled.
func (s *ReportService) CancellationRate(ctx context.Context, schema string) (CancellationRateDTO, error) {
	if _, err := s.db.ExecContext(ctx, fmt.Sprintf("SET search_path=%s", schema)); err != nil {
		return CancellationRateDTO{}, err
	}
	var resp CancellationRateDTO
	query := `
        SELECT SUM(CASE WHEN status = 'canceled' THEN 1 ELSE 0 END)::float / COUNT(*) AS rate
        FROM orders`
	if err := s.db.NewRaw(query).Scan(ctx, &resp); err != nil {
		return CancellationRateDTO{}, err
	}
	return resp, nil
}

// TODO: implement CancellationReasons (requires cancel_reason field)

// CurrentQueueLengthDTO holds current size of order queue.
type CurrentQueueLengthDTO struct {
	Length int `bun:"length"`
}

// CurrentQueueLength returns the number of active queued items (left_at IS NULL).
func (s *ReportService) CurrentQueueLength(ctx context.Context, schema string) (CurrentQueueLengthDTO, error) {
	if _, err := s.db.ExecContext(ctx, fmt.Sprintf("SET search_path=%s", schema)); err != nil {
		return CurrentQueueLengthDTO{}, err
	}
	var resp CurrentQueueLengthDTO
	query := `
        SELECT COUNT(*) AS length
        FROM order_queues
        WHERE left_at IS NULL`
	if err := s.db.NewRaw(query).Scan(ctx, &resp); err != nil {
		return CurrentQueueLengthDTO{}, err
	}
	return resp, nil
}

// AvgDeliveryTimeDTO holds average delivery time per driver.
type AvgDeliveryTimeDTO struct {
	DriverID   string  `bun:"driver_id"`
	AvgSeconds float64 `bun:"avg_seconds"`
}

// AvgDeliveryTimeByDriver returns average time from shipped to delivered per driver.
func (s *ReportService) AvgDeliveryTimeByDriver(ctx context.Context, schema string) ([]AvgDeliveryTimeDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []AvgDeliveryTimeDTO
	query := `
        SELECT driver_id::text AS driver_id,
            AVG(EXTRACT(EPOCH FROM (delivered_at - shipped_at))) AS avg_seconds
        FROM order_deliveries
        WHERE delivered_at IS NOT NULL AND shipped_at IS NOT NULL
        GROUP BY driver_id`
	if err := s.db.NewRaw(query).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeliveriesCountByDriverDTO holds count of deliveries per driver.
type DeliveriesCountByDriverDTO struct {
	DriverID string `bun:"driver_id"`
	Count    int    `bun:"count"`
}

// DeliveriesPerDriver returns number of deliveries made by each driver.
func (s *ReportService) DeliveriesPerDriver(ctx context.Context, schema string) ([]DeliveriesCountByDriverDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []DeliveriesCountByDriverDTO
	query := `
        SELECT driver_id::text AS driver_id, COUNT(*) AS count
        FROM order_deliveries
        GROUP BY driver_id`
	if err := s.db.NewRaw(query).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// OrdersPerTableDTO holds count of orders per table (utilization proxy).
type OrdersPerTableDTO struct {
	TableID string `bun:"table_id"`
	Count   int    `bun:"count"`
}

// OrdersPerTable returns number of orders per table.
func (s *ReportService) OrdersPerTable(ctx context.Context, schema string) ([]OrdersPerTableDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []OrdersPerTableDTO
	query := `
        SELECT table_id::text AS table_id, COUNT(*) AS count
        FROM order_tables
        GROUP BY table_id`
	if err := s.db.NewRaw(query).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SalesByShiftDTO holds sales aggregated by shift.
type SalesByShiftDTO struct {
	ShiftID string          `bun:"shift_id"`
	Total   decimal.Decimal `bun:"total"`
}

// SalesByShift returns total sales per shift.
func (s *ReportService) SalesByShift(ctx context.Context, schema string, start, end time.Time) ([]SalesByShiftDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []SalesByShiftDTO
	query := `
        SELECT shift_id::text AS shift_id, SUM(total_payable) AS total
        FROM orders
        WHERE created_at BETWEEN ? AND ?
        GROUP BY shift_id`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// PaymentsByMethodDTO holds sum of payments grouped by method.
type PaymentsByMethodDTO struct {
	Method string          `bun:"method"`
	Total  decimal.Decimal `bun:"total"`
}

// PaymentsByMethod returns sum of order payments per method.
func (s *ReportService) PaymentsByMethod(ctx context.Context, schema string, start, end time.Time) ([]PaymentsByMethodDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []PaymentsByMethodDTO
	query := `
        SELECT method, SUM(total_paid) AS total
        FROM order_payments
        WHERE paid_at BETWEEN ? AND ?
        GROUP BY method`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// EmployeePaymentsDTO holds sum of employee payments.
type EmployeePaymentsDTO struct {
	EmployeeID string          `bun:"employee_id"`
	Total      decimal.Decimal `bun:"total"`
}

// EmployeePaymentsReport returns sum of employee payments by employee.
func (s *ReportService) EmployeePaymentsReport(ctx context.Context, schema string, start, end time.Time) ([]EmployeePaymentsDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []EmployeePaymentsDTO
	query := `
        SELECT employee_id::text AS employee_id, SUM(amount) AS total
        FROM employee_payments
        WHERE pay_date BETWEEN ? AND ?
        GROUP BY employee_id`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: implement remaining methods 26–35...
