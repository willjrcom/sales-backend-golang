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

// NewReportRepository creates a new ReportService using the given DB connection.
func NewReportRepository(db *bun.DB) *ReportService {
	return &ReportService{db: db}
}

// SalesByDayDTO holds total payable grouped by day.
type SalesByDayDTO struct {
	Day   string          `bun:"day"`
	Total decimal.Decimal `bun:"total"`
}

// SalesTotalByDay returns total sales (sum of total_payable) per day in the given period.
func (s *ReportService) SalesTotalByDay(ctx context.Context, start, end time.Time) ([]SalesByDayDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []SalesByDayDTO
	query := `
        SELECT TO_CHAR(created_at, 'DD/MM') AS day, SUM(total_payable) AS total
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
func (s *ReportService) RevenueCumulativeByMonth(ctx context.Context, start, end time.Time) ([]CumulativeRevenueDTO, error) {
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
func (s *ReportService) SalesByHour(ctx context.Context, day time.Time) ([]HourlySalesDTO, error) {
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
func (s *ReportService) SalesByChannel(ctx context.Context, start, end time.Time) ([]ChannelSalesDTO, error) {
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
func (s *ReportService) AvgTicketByDay(ctx context.Context, start, end time.Time) ([]AvgTicketDTO, error) {
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
func (s *ReportService) AvgTicketByChannel(ctx context.Context, start, end time.Time) ([]AvgTicketDTO, error) {
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
	Day      string  `bun:"day"`
	Quantity float64 `bun:"quantity"`
}

// ProductsSoldByDay returns sum of quantity_items per day.
func (s *ReportService) ProductsSoldByDay(ctx context.Context, start, end time.Time) ([]ProductsSoldByDayDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []ProductsSoldByDayDTO
	query := `
		SELECT TO_CHAR(created_at, 'DD/MM') AS day, SUM(quantity_items) AS quantity
		FROM orders
		WHERE created_at BETWEEN ? AND ?
		GROUP BY day
		ORDER BY day
	`

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
func (s *ReportService) TopProducts(ctx context.Context, start, end time.Time) ([]TopProductsDTO, error) {
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
func (s *ReportService) SalesByCategory(ctx context.Context, start, end time.Time) ([]SalesByCategoryDTO, error) {
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
	if err := s.db.NewRaw(query, start.Format("2005-01-02"), end.Format("2005-01-02")).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ClientsRegisteredDTO holds count of new clients by day.
type ClientsRegisteredDTO struct {
	Day   string `bun:"day"`
	Count int    `bun:"count"`
}

// ClientsRegisteredByDay returns the number of clients registered per day.
func (s *ReportService) ClientsRegisteredByDay(ctx context.Context, start, end time.Time) ([]ClientsRegisteredDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []ClientsRegisteredDTO
	query := `
        SELECT TO_CHAR(created_at, 'DD/MM') AS day, COUNT(*) AS count
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
func (s *ReportService) NewVsRecurringClients(ctx context.Context, start, end time.Time) ([]NewVsRecurringDTO, error) {
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
func (s *ReportService) OrdersByStatus(ctx context.Context) ([]OrdersByStatusDTO, error) {
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
	ProcessRuleName string  `bun:"process_rule_name"`
	AvgSeconds      float64 `bun:"avg_seconds"`
}

// AvgProcessStepDurationByRule returns average time per process rule.
func (s *ReportService) AvgProcessStepDurationByRule(ctx context.Context) ([]AvgProcessStepDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []AvgProcessStepDTO
	query := `
        SELECT pr.name AS process_rule_name,
			AVG(EXTRACT(EPOCH FROM (op.finished_at - op.started_at))) AS avg_seconds
		FROM order_processes op
		JOIN process_rules pr ON pr.id = op.process_rule_id
		WHERE op.finished_at IS NOT NULL
		AND op.started_at IS NOT NULL
		GROUP BY pr.name`
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
func (s *ReportService) CancellationRate(ctx context.Context) (CancellationRateDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
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
func (s *ReportService) CurrentQueueLength(ctx context.Context) (CurrentQueueLengthDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
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

// AvgQueueDurationDTO holds average queue duration in seconds.
type AvgQueueDurationDTO struct {
	AvgSeconds float64 `bun:"avg_seconds"`
}

// AvgQueueDuration returns the average duration (in seconds) of all process queues.
func (s *ReportService) AvgQueueDuration(ctx context.Context) (AvgQueueDurationDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return AvgQueueDurationDTO{}, err
	}
	var resp AvgQueueDurationDTO
	// duration is stored as bigint nanoseconds; convert to seconds
	query := `
		SELECT AVG(duration) / 1000000000.0 AS avg_seconds
		FROM order_queues`
	if err := s.db.NewRaw(query).Scan(ctx, &resp); err != nil {
		return AvgQueueDurationDTO{}, err
	}
	return resp, nil
}

// AvgProcessByProductDTO holds average process duration per product.
type AvgProcessByProductDTO struct {
	ProductID   string  `bun:"product_id"`
	ProductName string  `bun:"product_name"`
	AvgSeconds  float64 `bun:"avg_seconds"`
}

// AvgProcessDurationByProduct returns average duration (seconds) of processes by product.
func (s *ReportService) AvgProcessDurationByProduct(ctx context.Context) ([]AvgProcessByProductDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}
	var resp []AvgProcessByProductDTO
	// p.duration is bigint nanoseconds; average and convert to seconds
	query := `
       SELECT prod.id::text AS product_id,
              prod.name AS product_name,
              AVG(p.duration) / 1000000000.0 AS avg_seconds
         FROM process_to_product_to_group_item t
         JOIN order_processes p ON p.id = t.process_id
         JOIN products prod ON prod.id = t.product_id
        GROUP BY prod.id, prod.name
        ORDER BY prod.name`
	if err := s.db.NewRaw(query).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// TotalQueueTimeByGroupItemDTO holds total queue time per group item in seconds.
type TotalQueueTimeByGroupItemDTO struct {
	Name         string  `bun:"name"`
	TotalSeconds float64 `bun:"total_seconds"`
}

// TotalQueueTimeByGroupItem returns total sum of queue durations per group item.
func (s *ReportService) TotalQueueTimeByGroupItem(ctx context.Context) ([]TotalQueueTimeByGroupItemDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}
	var resp []TotalQueueTimeByGroupItemDTO
	// duration is bigint nanoseconds; sum and convert to seconds
	query := `
       SELECT pr.name::text AS name,
              SUM(duration) / 1000000000.0 AS total_seconds
         FROM order_queues oq
		 JOIN process_rules pr ON pr.id = oq.process_rule_id
        GROUP BY pr.name`
	if err := s.db.NewRaw(query).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AvgDeliveryTimeDTO holds average delivery time per driver.
type AvgDeliveryTimeDTO struct {
	DriverID   string  `bun:"driver_id"`
	AvgSeconds float64 `bun:"avg_seconds"`
}

// AvgDeliveryTimeByDriver returns average time from shipped to delivered per driver.
func (s *ReportService) AvgDeliveryTimeByDriver(ctx context.Context) ([]AvgDeliveryTimeDTO, error) {
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
func (s *ReportService) DeliveriesPerDriver(ctx context.Context) ([]DeliveriesCountByDriverDTO, error) {
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

// TopTablesDTO holds the table usage count.
type TopTablesDTO struct {
	TableID string `bun:"table_id"`
	Count   int    `bun:"count"`
}

// TopTables returns the top 10 most used tables (by order count) in a period.
func (s *ReportService) TopTables(ctx context.Context, start, end time.Time) ([]TopTablesDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}
	var resp []TopTablesDTO
	query := `
       SELECT ot.table_id::text AS table_id,
              COUNT(*) AS count
         FROM order_tables ot
         JOIN orders o ON o.id = ot.order_id
        WHERE o.created_at BETWEEN ? AND ?
        GROUP BY ot.table_id
        ORDER BY count DESC
        LIMIT 10`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
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
func (s *ReportService) OrdersPerTable(ctx context.Context) ([]OrdersPerTableDTO, error) {
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
func (s *ReportService) SalesByShift(ctx context.Context, start, end time.Time) ([]SalesByShiftDTO, error) {
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
func (s *ReportService) PaymentsByMethod(ctx context.Context, start, end time.Time) ([]PaymentsByMethodDTO, error) {
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
func (s *ReportService) EmployeePaymentsReport(ctx context.Context, start, end time.Time) ([]EmployeePaymentsDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []EmployeePaymentsDTO
	query := `
        SELECT employee_id::text AS employee_id, SUM(amount) AS total
        FROM employee_payments
        WHERE payment_date BETWEEN ? AND ?
        GROUP BY employee_id`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SalesByPlaceDTO holds total sales per place.
type SalesByPlaceDTO struct {
	Place string          `bun:"place"`
	Total decimal.Decimal `bun:"total"`
}

// SalesByPlace returns total sales per place within the period.
func (s *ReportService) SalesByPlace(ctx context.Context, start, end time.Time) ([]SalesByPlaceDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []SalesByPlaceDTO
	query := `
        SELECT pl.name AS place, SUM(o.total_payable) AS total
        FROM orders o
        JOIN order_tables ot ON ot.order_id = o.id
        JOIN place_to_tables pt ON pt.table_id = ot.table_id
        JOIN places pl ON pl.id = pt.place_id
        WHERE o.created_at BETWEEN ? AND ?
        GROUP BY pl.name`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SalesBySizeDTO holds total quantity sold per product size.
type SalesBySizeDTO struct {
	Size     string  `bun:"size"`
	Quantity float64 `bun:"quantity"`
}

// SalesBySize returns total quantity sold grouped by product size.
func (s *ReportService) SalesBySize(ctx context.Context, start, end time.Time) ([]SalesBySizeDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []SalesBySizeDTO
	query := `
        SELECT i.size AS size, SUM(i.quantity) AS quantity
        FROM order_items i
        JOIN order_group_items g ON g.id = i.group_item_id
        JOIN orders o ON o.id = g.order_id
        WHERE o.created_at BETWEEN ? AND ?
        GROUP BY i.size`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AdditionalItemsDTO holds total quantity of additional items sold.
type AdditionalItemsDTO struct {
	Name     string  `bun:"name"`
	Quantity float64 `bun:"quantity"`
}

// AdditionalItemsSold returns total quantity of additional items sold.
func (s *ReportService) AdditionalItemsSold(ctx context.Context, start, end time.Time) ([]AdditionalItemsDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []AdditionalItemsDTO
	query := `
        SELECT i.name AS name, SUM(i.quantity) AS quantity
        FROM order_items i
        JOIN order_group_items g ON g.id = i.group_item_id
        JOIN orders o ON o.id = g.order_id
        WHERE o.created_at BETWEEN ? AND ? AND i.is_additional = TRUE
        GROUP BY i.name`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AvgPickupTimeDTO holds average waiting time for pickups in seconds.
type AvgPickupTimeDTO struct {
	AvgSeconds float64 `bun:"avg_seconds"`
}

// AvgPickupTime returns average time between pending and ready for pickups.
func (s *ReportService) AvgPickupTime(ctx context.Context, start, end time.Time) (AvgPickupTimeDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return AvgPickupTimeDTO{}, err
	}

	var resp AvgPickupTimeDTO
	query := `
        SELECT AVG(EXTRACT(EPOCH FROM (ready_at - pending_at))) AS avg_seconds
        FROM order_pickups p
        JOIN orders o ON o.id = p.order_id
        WHERE p.ready_at IS NOT NULL AND p.pending_at IS NOT NULL AND o.created_at BETWEEN ? AND ?`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return AvgPickupTimeDTO{}, err
	}
	return resp, nil
}

// GroupItemsStatusDTO holds count of group items by status.
type GroupItemsStatusDTO struct {
	Status string `bun:"status"`
	Count  int    `bun:"count"`
}

// GroupItemsByStatus returns count of group items per status.
func (s *ReportService) GroupItemsByStatus(ctx context.Context, start, end time.Time) ([]GroupItemsStatusDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []GroupItemsStatusDTO
	query := `
        SELECT g.status AS status, COUNT(*) AS count
        FROM order_group_items g
        JOIN orders o ON o.id = g.order_id
        WHERE o.created_at BETWEEN ? AND ?
        GROUP BY g.status`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeliveriesByCepDTO holds count of deliveries per CEP.
type DeliveriesByCepDTO struct {
	Cep   string `bun:"cep"`
	Count int    `bun:"count"`
}

// DeliveriesByCep returns number of deliveries per ZIP code.
func (s *ReportService) DeliveriesByCep(ctx context.Context, start, end time.Time) ([]DeliveriesByCepDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []DeliveriesByCepDTO
	query := `
        SELECT a.cep AS cep, COUNT(*) AS count
        FROM order_deliveries d
        JOIN addresses a ON a.id = d.address_id
        JOIN orders o ON o.id = d.order_id
        WHERE d.delivered_at BETWEEN ? AND ?
        GROUP BY a.cep`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ProcessedByRuleDTO holds count of orders processed per process rule.
type ProcessedByRuleDTO struct {
	RuleID string `bun:"rule_id"`
	Count  int    `bun:"count"`
}

// ProcessedCountByRule returns number of processed items per rule.
func (s *ReportService) ProcessedCountByRule(ctx context.Context, start, end time.Time) ([]ProcessedByRuleDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []ProcessedByRuleDTO
	query := `
        SELECT pr.process_rule_id::text AS rule_id, COUNT(*) AS count
        FROM order_processes pr
        JOIN order_group_items g ON g.id = pr.group_item_id
        JOIN orders o ON o.id = g.order_id
        WHERE pr.finished_at IS NOT NULL AND pr.started_at IS NOT NULL
        	AND o.created_at BETWEEN ? AND ?
        GROUP BY pr.process_rule_id`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DailySalesDTO holds summary of sales for a specific day.
type DailySalesDTO struct {
	TotalOrders int             `bun:"total_orders"`
	TotalSales  decimal.Decimal `bun:"total_sales"`
}

// DailySales returns summary metrics for the given day.
func (s *ReportService) DailySales(ctx context.Context, day time.Time) (DailySalesDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return DailySalesDTO{}, err
	}

	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)
	var resp DailySalesDTO
	query := `
        SELECT COUNT(*) AS total_orders, SUM(total_payable) AS total_sales
        FROM orders
        WHERE created_at >= ? AND created_at < ?`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return DailySalesDTO{}, err
	}
	return resp, nil
}

// ProductProfitabilityDTO holds profitability metrics per product.
type ProductProfitabilityDTO struct {
	ProductID    string           `bun:"product_id"`
	ProductName  string           `bun:"product_name"`
	TotalSold    *decimal.Decimal `bun:"total_sold"`
	TotalCost    *decimal.Decimal `bun:"total_cost"`
	TotalRevenue *decimal.Decimal `bun:"total_revenue"`
	Profit       *decimal.Decimal `bun:"profit"`
	ProfitMargin *decimal.Decimal `bun:"profit_margin"`
}

// ProductProfitability returns profitability analysis per product.
func (s *ReportService) ProductProfitability(ctx context.Context, start, end time.Time) ([]ProductProfitabilityDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []ProductProfitabilityDTO
	query := `
        SELECT 
            p.id::text AS product_id,
            p.name AS product_name,
            COALESCE(SUM(i.quantity), 0) AS total_sold,
            COALESCE(SUM(i.quantity * p.cost), 0) AS total_cost,
            COALESCE(SUM(i.quantity * i.price), 0) AS total_revenue,
            COALESCE(SUM(i.quantity * i.price) - SUM(i.quantity * p.cost), 0) AS profit,
            CASE 
                WHEN COALESCE(SUM(i.quantity * i.price), 0) > 0 
                THEN ((COALESCE(SUM(i.quantity * i.price), 0) - COALESCE(SUM(i.quantity * p.cost), 0)) / COALESCE(SUM(i.quantity * i.price), 0)) * 100
                ELSE 0 
            END AS profit_margin
        FROM order_items i
        JOIN order_group_items g ON g.id = i.group_item_id
        JOIN orders o ON o.id = g.order_id
        JOIN products p ON p.id = i.product_id
        WHERE o.created_at BETWEEN ? AND ?
        GROUP BY p.id, p.name
        ORDER BY profit DESC`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CategoryProfitabilityDTO holds profitability metrics per category.
type CategoryProfitabilityDTO struct {
	CategoryName string           `bun:"category_name"`
	TotalSold    *decimal.Decimal `bun:"total_sold"`
	TotalCost    *decimal.Decimal `bun:"total_cost"`
	TotalRevenue *decimal.Decimal `bun:"total_revenue"`
	Profit       *decimal.Decimal `bun:"profit"`
	ProfitMargin *decimal.Decimal `bun:"profit_margin"`
}

// CategoryProfitability returns profitability analysis per product category.
func (s *ReportService) CategoryProfitability(ctx context.Context, start, end time.Time) ([]CategoryProfitabilityDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	var resp []CategoryProfitabilityDTO
	query := `
        SELECT 
            pc.name AS category_name,
            COALESCE(SUM(i.quantity), 0) AS total_sold,
            COALESCE(SUM(i.quantity * p.cost), 0) AS total_cost,
            COALESCE(SUM(i.quantity * i.price), 0) AS total_revenue,
            COALESCE(SUM(i.quantity * i.price) - SUM(i.quantity * p.cost), 0) AS profit,
            CASE 
                WHEN COALESCE(SUM(i.quantity * i.price), 0) > 0 
                THEN ((COALESCE(SUM(i.quantity * i.price), 0) - COALESCE(SUM(i.quantity * p.cost), 0)) / COALESCE(SUM(i.quantity * i.price), 0)) * 100
                ELSE 0 
            END AS profit_margin
        FROM order_items i
        JOIN order_group_items g ON g.id = i.group_item_id
        JOIN orders o ON o.id = g.order_id
        JOIN products p ON p.id = i.product_id
        JOIN product_categories pc ON pc.id = p.category_id
        WHERE o.created_at BETWEEN ? AND ?
        GROUP BY pc.name
        ORDER BY profit DESC`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// LowProfitProductsDTO holds products with low profit margin.
type LowProfitProductsDTO struct {
	ProductID    string           `bun:"product_id"`
	ProductName  string           `bun:"product_name"`
	Price        *decimal.Decimal `bun:"price"`
	Cost         *decimal.Decimal `bun:"cost"`
	ProfitMargin *decimal.Decimal `bun:"profit_margin"`
}

// LowProfitProducts returns products with profit margin below threshold.
func (s *ReportService) LowProfitProducts(ctx context.Context, minMargin float64) ([]LowProfitProductsDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	// Debug: Primeiro vamos verificar se há produtos cadastrados
	var debugCount int
	debugQuery := `SELECT COUNT(*) FROM products WHERE is_available = true`
	if err := s.db.NewRaw(debugQuery).Scan(ctx, &debugCount); err != nil {
		return nil, err
	}

	// Debug: Verificar produtos com custo
	var productsWithCost int
	costQuery := `SELECT COUNT(*) FROM products WHERE is_available = true AND cost > 0`
	if err := s.db.NewRaw(costQuery).Scan(ctx, &productsWithCost); err != nil {
		return nil, err
	}

	// Debug: Verificar produtos sem custo
	var productsWithoutCost int
	noCostQuery := `SELECT COUNT(*) FROM products WHERE is_available = true AND (cost = 0 OR cost IS NULL)`
	if err := s.db.NewRaw(noCostQuery).Scan(ctx, &productsWithoutCost); err != nil {
		return nil, err
	}

	fmt.Printf("Debug - Total produtos disponíveis: %d\n", debugCount)
	fmt.Printf("Debug - Produtos com custo > 0: %d\n", productsWithCost)
	fmt.Printf("Debug - Produtos sem custo: %d\n", productsWithoutCost)
	fmt.Printf("Debug - Margem mínima configurada: %.2f%%\n", minMargin)

	var resp []LowProfitProductsDTO
	query := `
        SELECT 
            p.id::text AS product_id,
            p.name AS product_name,
            p.price,
            p.cost,
            CASE 
                WHEN p.price > 0 AND p.cost > 0
                THEN ((p.price - p.cost) / p.price) * 100
                WHEN p.price > 0 AND (p.cost = 0 OR p.cost IS NULL)
                THEN 100
                ELSE 0 
            END AS profit_margin
        FROM products p
        WHERE p.price > 0 
        AND p.is_available = true
        AND (
            -- Produtos com margem baixa
            (p.cost > 0 AND ((p.price - p.cost) / p.price) * 100 < ?)
            OR 
            -- Produtos sem custo definido (custo = 0 ou NULL)
            (p.cost = 0 OR p.cost IS NULL)
        )
        ORDER BY profit_margin ASC`
	if err := s.db.NewRaw(query, minMargin).Scan(ctx, &resp); err != nil {
		return nil, err
	}

	fmt.Printf("Debug - Produtos retornados: %d\n", len(resp))
	return resp, nil
}

// DebugProducts returns debug information about products
func (s *ReportService) DebugProducts(ctx context.Context) (map[string]interface{}, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return nil, err
	}

	// Total produtos
	var totalProducts int
	totalQuery := `SELECT COUNT(*) FROM products`
	if err := s.db.NewRaw(totalQuery).Scan(ctx, &totalProducts); err != nil {
		return nil, err
	}

	// Produtos disponíveis
	var availableProducts int
	availableQuery := `SELECT COUNT(*) FROM products WHERE is_available = true`
	if err := s.db.NewRaw(availableQuery).Scan(ctx, &availableProducts); err != nil {
		return nil, err
	}

	// Produtos com custo
	var productsWithCost int
	costQuery := `SELECT COUNT(*) FROM products WHERE is_available = true AND cost > 0`
	if err := s.db.NewRaw(costQuery).Scan(ctx, &productsWithCost); err != nil {
		return nil, err
	}

	// Produtos sem custo
	var productsWithoutCost int
	noCostQuery := `SELECT COUNT(*) FROM products WHERE is_available = true AND (cost = 0 OR cost IS NULL)`
	if err := s.db.NewRaw(noCostQuery).Scan(ctx, &productsWithoutCost); err != nil {
		return nil, err
	}

	// Produtos com preço
	var productsWithPrice int
	priceQuery := `SELECT COUNT(*) FROM products WHERE is_available = true AND price > 0`
	if err := s.db.NewRaw(priceQuery).Scan(ctx, &productsWithPrice); err != nil {
		return nil, err
	}

	// Lista de produtos para debug
	var products []map[string]interface{}
	productsQuery := `
		SELECT 
			id::text as product_id,
			name as product_name,
			price,
			cost,
			is_available,
			CASE 
				WHEN price > 0 AND cost > 0
				THEN ((price - cost) / price) * 100
				WHEN price > 0 AND (cost = 0 OR cost IS NULL)
				THEN 100
				ELSE 0 
			END AS profit_margin
		FROM products 
		WHERE is_available = true 
		ORDER BY name`

	if err := s.db.NewRaw(productsQuery).Scan(ctx, &products); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_products":        totalProducts,
		"available_products":    availableProducts,
		"products_with_cost":    productsWithCost,
		"products_without_cost": productsWithoutCost,
		"products_with_price":   productsWithPrice,
		"products":              products,
	}, nil
}

// OverallProfitabilityDTO holds overall profitability metrics.
type OverallProfitabilityDTO struct {
	TotalRevenue *decimal.Decimal `bun:"total_revenue"`
	TotalCost    *decimal.Decimal `bun:"total_cost"`
	TotalProfit  *decimal.Decimal `bun:"total_profit"`
	ProfitMargin *decimal.Decimal `bun:"profit_margin"`
}

// OverallProfitability returns overall profitability metrics for the period.
func (s *ReportService) OverallProfitability(ctx context.Context, start, end time.Time) (OverallProfitabilityDTO, error) {
	if err := database.ChangeSchema(ctx, s.db); err != nil {
		return OverallProfitabilityDTO{}, err
	}

	var resp OverallProfitabilityDTO
	query := `
        SELECT 
            COALESCE(SUM(i.quantity * i.price), 0) AS total_revenue,
            COALESCE(SUM(i.quantity * p.cost), 0) AS total_cost,
            COALESCE(SUM(i.quantity * i.price) - SUM(i.quantity * p.cost), 0) AS total_profit,
            CASE 
                WHEN COALESCE(SUM(i.quantity * i.price), 0) > 0 
                THEN ((COALESCE(SUM(i.quantity * i.price), 0) - COALESCE(SUM(i.quantity * p.cost), 0)) / COALESCE(SUM(i.quantity * i.price), 0)) * 100
                ELSE 0 
            END AS profit_margin
        FROM order_items i
        JOIN order_group_items g ON g.id = i.group_item_id
        JOIN orders o ON o.id = g.order_id
        JOIN products p ON p.id = i.product_id
        WHERE o.created_at BETWEEN ? AND ?`
	if err := s.db.NewRaw(query, start, end).Scan(ctx, &resp); err != nil {
		return OverallProfitabilityDTO{}, err
	}
	return resp, nil
}
