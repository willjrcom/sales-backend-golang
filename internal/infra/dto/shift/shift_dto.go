package shiftdto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
	employeedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/employee"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

type ShiftDTO struct {
	ID uuid.UUID `json:"id"`
	ShiftTimeLogs
	ShiftCommonAttributes
}

type ShiftCommonAttributes struct {
	CurrentOrderNumber     int                        `json:"current_order_number"`
	Orders                 []orderdto.OrderDTO        `json:"orders"`
	Redeems                []RedeemDTO                `json:"redeems"`
	StartChange            decimal.Decimal            `json:"start_change"`
	EndChange              *decimal.Decimal           `json:"end_change"`
	AttendantID            *uuid.UUID                 `json:"attendant_id"`
	Attendant              *employeedto.EmployeeDTO   `json:"attendant"`
	TotalOrdersFinished    int                        `json:"total_orders_finished"`
	TotalOrdersCancelled   int                        `json:"total_orders_cancelled"`
	TotalSales             decimal.Decimal            `json:"total_sales"`
	SalesByCategory        map[string]decimal.Decimal `json:"sales_by_category"`
	ProductsSoldByCategory map[string]float64         `json:"products_sold_by_category"`
	TotalItemsSold         float64                    `json:"total_items_sold"`
	AverageOrderValue      decimal.Decimal            `json:"average_order_value"`
	Payments               []orderdto.PaymentOrderDTO `json:"payments"`
	DeliveryDrivers        []DeliveryDriverTaxDTO     `json:"delivery_drivers_tax"`
	// Campos de analytics de produção
	OrderProcessAnalytics  map[string]*OrderProcessAnalyticsDTO `json:"order_process_analytics"`
	QueueAnalytics         map[string]*QueueAnalyticsDTO        `json:"queue_analytics"`
	TotalProcesses         int                                  `json:"total_processes"`
	TotalQueues            int                                  `json:"total_queues"`
	AverageProcessTime     int64                                `json:"average_process_time"` // em segundos
	AverageQueueTime       int64                                `json:"average_queue_time"`   // em segundos
	ProcessEfficiencyScore decimal.Decimal                      `json:"process_efficiency_score"`
}

type ShiftTimeLogs struct {
	OpenedAt *time.Time `json:"opened_at"`
	ClosedAt *time.Time `json:"closed_at"`
}

// OrderProcessAnalyticsDTO representa as métricas de uma regra de processo
type OrderProcessAnalyticsDTO struct {
	ProcessRuleID       uuid.UUID                             `json:"process_rule_id"`
	ProcessRuleName     string                                `json:"process_rule_name"`
	TotalProcesses      int                                   `json:"total_processes"`
	CompletedProcesses  int                                   `json:"completed_processes"`
	CancelledProcesses  int                                   `json:"cancelled_processes"`
	AverageProcessTime  int64                                 `json:"average_process_time"` // em segundos
	TotalProcessTime    int64                                 `json:"total_process_time"`   // em segundos
	TotalPausedCount    int                                   `json:"total_paused_count"`
	EfficiencyScore     decimal.Decimal                       `json:"efficiency_score"`
	CategoriesProcessed map[string]*CategoryProcessMetricsDTO `json:"categories_processed"`
	EmployeePerformance map[string]*EmployeeProcessMetricsDTO `json:"employee_performance"`
}

// QueueAnalyticsDTO representa as métricas de fila
type QueueAnalyticsDTO struct {
	ProcessRuleID    uuid.UUID `json:"process_rule_id"`
	ProcessRuleName  string    `json:"process_rule_name"`
	TotalQueues      int       `json:"total_queues"`
	CompletedQueues  int       `json:"completed_queues"`
	AverageQueueTime int64     `json:"average_queue_time"` // em segundos
	TotalQueueTime   int64     `json:"total_queue_time"`   // em segundos
}

// CategoryProcessMetricsDTO representa métricas por categoria
type CategoryProcessMetricsDTO struct {
	CategoryID         uuid.UUID `json:"category_id"`
	CategoryName       string    `json:"category_name"`
	TotalProcesses     int       `json:"total_processes"`
	AverageProcessTime int64     `json:"average_process_time"` // em segundos
}

// EmployeeProcessMetricsDTO representa métricas por funcionário
type EmployeeProcessMetricsDTO struct {
	EmployeeID         uuid.UUID       `json:"employee_id"`
	EmployeeName       string          `json:"employee_name"`
	TotalProcesses     int             `json:"total_processes"`
	CompletedProcesses int             `json:"completed_processes"`
	AverageProcessTime int64           `json:"average_process_time"` // em segundos
	EfficiencyScore    decimal.Decimal `json:"efficiency_score"`
}

func (s *ShiftDTO) FromDomain(shift *shiftentity.Shift) {
	if shift == nil {
		return
	}
	*s = ShiftDTO{
		ID: shift.ID,
		ShiftTimeLogs: ShiftTimeLogs{
			OpenedAt: shift.OpenedAt,
			ClosedAt: shift.ClosedAt,
		},
		ShiftCommonAttributes: ShiftCommonAttributes{
			CurrentOrderNumber:     shift.CurrentOrderNumber,
			Orders:                 []orderdto.OrderDTO{},
			Redeems:                []RedeemDTO{},
			StartChange:            shift.StartChange,
			EndChange:              shift.EndChange,
			AttendantID:            shift.AttendantID,
			Attendant:              &employeedto.EmployeeDTO{},
			TotalOrdersFinished:    shift.TotalOrdersFinished,
			TotalOrdersCancelled:   shift.TotalOrdersCancelled,
			TotalSales:             shift.TotalSales,
			SalesByCategory:        shift.SalesByCategory,
			ProductsSoldByCategory: shift.ProductsSoldByCategory,
			TotalItemsSold:         shift.TotalItemsSold,
			AverageOrderValue:      shift.AverageOrderValue,
			Payments:               []orderdto.PaymentOrderDTO{},
			DeliveryDrivers:        []DeliveryDriverTaxDTO{},
			// Analytics de produção
			OrderProcessAnalytics:  make(map[string]*OrderProcessAnalyticsDTO),
			QueueAnalytics:         make(map[string]*QueueAnalyticsDTO),
			TotalProcesses:         shift.TotalProcesses,
			TotalQueues:            shift.TotalQueues,
			AverageProcessTime:     int64(shift.AverageProcessTime.Seconds()),
			AverageQueueTime:       int64(shift.AverageQueueTime.Seconds()),
			ProcessEfficiencyScore: shift.ProcessEfficiencyScore,
		},
	}

	for _, order := range shift.Orders {
		o := orderdto.OrderDTO{}
		o.FromDomain(&order)
		s.Orders = append(s.Orders, o)
	}

	for _, redeem := range shift.Redeems {
		r := RedeemDTO{}
		r.FromDomain(&redeem)
		s.Redeems = append(s.Redeems, r)
	}

	for _, payment := range shift.Payments {
		p := orderdto.PaymentOrderDTO{}
		p.FromDomain(&payment)
		s.Payments = append(s.Payments, p)
	}

	for _, tax := range shift.DeliveryDrivers {
		t := DeliveryDriverTaxDTO{}
		t.FromDomain(&tax)
		s.DeliveryDrivers = append(s.DeliveryDrivers, t)
	}

	// Converte analytics de produção
	for _, analytics := range shift.OrderProcessAnalytics {
		s.OrderProcessAnalytics[analytics.ProcessRuleName] = &OrderProcessAnalyticsDTO{
			ProcessRuleID:       analytics.ProcessRuleID,
			ProcessRuleName:     analytics.ProcessRuleName,
			TotalProcesses:      analytics.TotalProcesses,
			CompletedProcesses:  analytics.CompletedProcesses,
			CancelledProcesses:  analytics.CancelledProcesses,
			AverageProcessTime:  int64(analytics.AverageProcessTime.Seconds()),
			TotalProcessTime:    int64(analytics.TotalProcessTime.Seconds()),
			TotalPausedCount:    analytics.TotalPausedCount,
			EfficiencyScore:     analytics.GetEfficiencyScore(5 * time.Minute), // 5 min esperado
			CategoriesProcessed: make(map[string]*CategoryProcessMetricsDTO),
			EmployeePerformance: make(map[string]*EmployeeProcessMetricsDTO),
		}

		// Converte categorias
		for categoryID, categoryMetrics := range analytics.CategoriesProcessed {
			s.OrderProcessAnalytics[analytics.ProcessRuleName].CategoriesProcessed[categoryMetrics.CategoryName] = &CategoryProcessMetricsDTO{
				CategoryID:         categoryID,
				CategoryName:       categoryMetrics.CategoryName,
				TotalProcesses:     categoryMetrics.TotalProcessed,
				AverageProcessTime: int64(categoryMetrics.AverageProcessTime.Seconds()),
			}
		}

		// Converte funcionários
		for employeeID, employeeMetrics := range analytics.EmployeePerformance {
			s.OrderProcessAnalytics[analytics.ProcessRuleName].EmployeePerformance[employeeMetrics.EmployeeName] = &EmployeeProcessMetricsDTO{
				EmployeeID:         employeeID,
				EmployeeName:       employeeMetrics.EmployeeName,
				TotalProcesses:     employeeMetrics.TotalProcessed,
				CompletedProcesses: employeeMetrics.TotalProcessed, // Assumindo que todos os processos processados foram completados
				AverageProcessTime: int64(employeeMetrics.AverageProcessTime.Seconds()),
				EfficiencyScore:    employeeMetrics.EfficiencyScore,
			}
		}
	}

	// Converte analytics de fila
	for _, queueAnalytics := range shift.QueueAnalytics {
		s.QueueAnalytics[queueAnalytics.ProcessRuleName] = &QueueAnalyticsDTO{
			ProcessRuleID:    queueAnalytics.ProcessRuleID,
			ProcessRuleName:  queueAnalytics.ProcessRuleName,
			TotalQueues:      queueAnalytics.TotalQueued,
			CompletedQueues:  queueAnalytics.TotalQueued, // Assumindo que todas as filas foram completadas
			AverageQueueTime: int64(queueAnalytics.AverageQueueTime.Seconds()),
			TotalQueueTime:   int64(queueAnalytics.TotalQueueTime.Seconds()),
		}
	}

	s.Attendant.FromDomain(shift.Attendant)

	if len(shift.Orders) == 0 {
		s.Orders = nil
	}
	if len(shift.Redeems) == 0 {
		s.Redeems = nil
	}
	if shift.Attendant == nil {
		s.Attendant = nil
	}
	if len(s.OrderProcessAnalytics) == 0 {
		s.OrderProcessAnalytics = nil
	}
	if len(s.QueueAnalytics) == 0 {
		s.QueueAnalytics = nil
	}
}
