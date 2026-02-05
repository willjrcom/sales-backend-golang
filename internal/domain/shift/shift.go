package shiftentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
)

type Shift struct {
	entity.Entity
	ShiftTimeLogs
	ShiftCommonAttributes
}

type ShiftCommonAttributes struct {
	CurrentOrderNumber     int
	Orders                 []orderentity.Order
	Redeems                []Redeem
	StartChange            decimal.Decimal
	EndChange              *decimal.Decimal
	AttendantID            *uuid.UUID
	Attendant              *employeeentity.Employee
	TotalOrdersFinished    int
	TotalOrdersCancelled   int
	TotalSales             decimal.Decimal
	SalesByCategory        map[string]decimal.Decimal
	ProductsSoldByCategory map[string]float64
	TotalItemsSold         float64         // soma de todas as quantidades de itens, para medir o "pulo de prato"
	AverageOrderValue      decimal.Decimal // TotalSales ÷ TotalOrders, para análise de ticket médio
	Payments               []orderentity.PaymentOrder
	DeliveryDrivers        []DeliveryDriverTax
	// Novos campos para analytics de produção
	OrderProcessAnalytics  map[uuid.UUID]*OrderProcessAnalytics // ProcessRuleID -> Analytics
	QueueAnalytics         map[uuid.UUID]*QueueAnalytics        // ProcessRuleID -> Queue Analytics
	TotalProcesses         int
	TotalQueues            int
	AverageProcessTime     time.Duration
	AverageQueueTime       time.Duration
	ProcessEfficiencyScore decimal.Decimal
}

type ShiftTimeLogs struct {
	OpenedAt *time.Time
	ClosedAt *time.Time
}

// NewShift creates a new shift with initial start change
func NewShift(startChange decimal.Decimal) *Shift {
	newEntity := entity.NewEntity()
	now := time.Now().UTC()
	newEntity.CreatedAt = now

	shift := &Shift{
		Entity: newEntity,
		ShiftCommonAttributes: ShiftCommonAttributes{
			CurrentOrderNumber:    0,
			StartChange:           startChange,
			OrderProcessAnalytics: make(map[uuid.UUID]*OrderProcessAnalytics),
			QueueAnalytics:        make(map[uuid.UUID]*QueueAnalytics),
		},
		ShiftTimeLogs: ShiftTimeLogs{
			OpenedAt: &now,
		},
	}

	return shift
}

// CloseShift closes the shift with final change
func (s *Shift) CloseShift(endChange decimal.Decimal) {
	now := time.Now().UTC()
	s.EndChange = &endChange
	s.ClosedAt = &now
}

func (s *Shift) Load(deliveryDrivers map[uuid.UUID]orderentity.DeliveryDriver, processes []orderprocessentity.OrderProcess, queues []*orderprocessentity.OrderQueue, processRules map[uuid.UUID]string, employees map[uuid.UUID]string) {
	// compute analytics for reporting
	s.TotalOrdersFinished = 0
	s.TotalOrdersCancelled = 0
	s.TotalSales = decimal.Zero

	// initialize maps
	s.SalesByCategory = make(map[string]decimal.Decimal)
	s.ProductsSoldByCategory = make(map[string]float64)
	s.TotalItemsSold = 0
	s.Payments = make([]orderentity.PaymentOrder, 0)
	s.DeliveryDrivers = make([]DeliveryDriverTax, 0)

	// aggregate orders data
	for _, o := range s.Orders {
		if o.Status == orderentity.OrderStatusCancelled {
			s.TotalOrdersCancelled++
			continue
		}

		if o.Status != orderentity.OrderStatusFinished {
			continue
		}

		s.TotalOrdersFinished++
		s.Payments = append(s.Payments, o.Payments...)

		if o.Delivery != nil && o.Delivery.DriverID != nil {
			deliveryDriver, ok := deliveryDrivers[*o.Delivery.DriverID]
			if ok {
				deliveryDriverTax := DeliveryDriverTax{
					DeliveryDriverID:   deliveryDriver.ID,
					DeliveryDriverName: deliveryDriver.Employee.User.Name,
					OrderNumber:        o.OrderNumber,
					DeliveryID:         o.Delivery.ID,
					DeliveryTax:        *o.Delivery.DeliveryTax,
				}

				s.DeliveryDrivers = append(s.DeliveryDrivers, deliveryDriverTax)
			}
		}

		// ensure totals are up to date
		s.TotalSales = s.TotalSales.Add(o.TotalPayable)
		for _, g := range o.GroupItems {
			cat := ""
			if g.Category != nil {
				cat = g.Category.Name
			}

			// sum revenue by category
			rev := g.TotalPrice
			if prev, ok := s.SalesByCategory[cat]; ok {
				s.SalesByCategory[cat] = prev.Add(rev)
			} else {
				s.SalesByCategory[cat] = rev
			}

			// sum quantities by category and total items
			qty := g.Quantity
			s.ProductsSoldByCategory[cat] += qty
			s.TotalItemsSold += qty
		}
	}

	s.AverageOrderValue = decimal.Zero
	// average order value
	if s.TotalOrdersFinished > 0 {
		s.AverageOrderValue = s.TotalSales.Div(decimal.NewFromInt(int64(s.TotalOrdersFinished)))
	}

	// Carrega métricas de produção se os dados estiverem disponíveis
	if processes != nil && queues != nil && processRules != nil && employees != nil {
		s.loadOrderProcessAnalytics(processes, queues, processRules, employees)
	}
}

// loadOrderProcessAnalytics carrega e calcula as métricas de processos de produção (método privado)
func (s *Shift) loadOrderProcessAnalytics(processes []orderprocessentity.OrderProcess, queues []*orderprocessentity.OrderQueue, processRules map[uuid.UUID]string, employees map[uuid.UUID]string) {
	s.TotalProcesses = 0
	s.TotalQueues = 0
	s.AverageProcessTime = 0
	s.AverageQueueTime = 0
	s.ProcessEfficiencyScore = decimal.Zero
	s.OrderProcessAnalytics = map[uuid.UUID]*OrderProcessAnalytics{}
	s.QueueAnalytics = map[uuid.UUID]*QueueAnalytics{}

	var totalProcessTime time.Duration
	var totalQueueTime time.Duration
	var completedProcesses int

	// Processa os processos
	for _, process := range processes {
		s.TotalProcesses++

		// Obtém o nome da regra de processo
		processRuleName := ""
		if name, exists := processRules[process.ProcessRuleID]; exists {
			processRuleName = name
		}

		// Obtém o nome do funcionário
		employeeName := ""
		if process.EmployeeID != nil {
			if name, exists := employees[*process.EmployeeID]; exists {
				employeeName = name
			}
		}

		// Inicializa analytics para esta regra de processo se não existir
		if _, exists := s.OrderProcessAnalytics[process.ProcessRuleID]; !exists {
			s.OrderProcessAnalytics[process.ProcessRuleID] = NewOrderProcessAnalytics(process.ProcessRuleID, processRuleName)
		}

		// Adiciona o processo às métricas
		s.OrderProcessAnalytics[process.ProcessRuleID].AddProcess(&process, employeeName)

		// Calcula métricas gerais
		if process.Status == orderprocessentity.ProcessStatusFinished {
			completedProcesses++
			totalProcessTime += process.Duration
		}
	}

	// Processa as filas
	for _, queue := range queues {
		s.TotalQueues++

		// Obtém o nome da regra de processo
		processRuleName := ""
		if name, exists := processRules[queue.ProcessRuleID]; exists {
			processRuleName = name
		}

		// Inicializa analytics para esta regra de processo se não existir
		if _, exists := s.QueueAnalytics[queue.ProcessRuleID]; !exists {
			s.QueueAnalytics[queue.ProcessRuleID] = NewQueueAnalytics(queue.ProcessRuleID, processRuleName)
		}

		// Adiciona a fila às métricas
		s.QueueAnalytics[queue.ProcessRuleID].AddQueue(queue)

		// Calcula métricas gerais
		if queue.LeftAt != nil {
			totalQueueTime += queue.Duration
		}
	}

	// Calcula médias gerais
	if completedProcesses > 0 {
		s.AverageProcessTime = time.Duration(int64(totalProcessTime) / int64(completedProcesses))
	}

	if s.TotalQueues > 0 {
		s.AverageQueueTime = time.Duration(int64(totalQueueTime) / int64(s.TotalQueues))
	}

	// Calcula score de eficiência geral (exemplo: tempo esperado de 5 minutos)
	expectedTime := 5 * time.Minute
	s.ProcessEfficiencyScore = s.calculateOverallEfficiencyScore(expectedTime)
}

// calculateOverallEfficiencyScore calcula o score de eficiência geral
func (s *Shift) calculateOverallEfficiencyScore(expectedTime time.Duration) decimal.Decimal {
	if s.AverageProcessTime == 0 || expectedTime == 0 {
		return decimal.Zero
	}

	// Score baseado na relação entre tempo esperado e tempo real
	ratio := float64(expectedTime) / float64(s.AverageProcessTime)
	return decimal.NewFromFloat(ratio * 100) // Score de 0-100
}

// GetProcessAnalyticsByRule retorna as métricas de uma regra de processo específica
func (s *Shift) GetProcessAnalyticsByRule(processRuleID uuid.UUID) *OrderProcessAnalytics {
	if analytics, exists := s.OrderProcessAnalytics[processRuleID]; exists {
		return analytics
	}
	return nil
}

// GetQueueAnalyticsByRule retorna as métricas de fila de uma regra de processo específica
func (s *Shift) GetQueueAnalyticsByRule(processRuleID uuid.UUID) *QueueAnalytics {
	if analytics, exists := s.QueueAnalytics[processRuleID]; exists {
		return analytics
	}
	return nil
}

// GetTopPerformingEmployees retorna os funcionários com melhor performance
func (s *Shift) GetTopPerformingEmployees(limit int) []EmployeeProcessMetrics {
	var allEmployees []EmployeeProcessMetrics

	// Coleta todos os funcionários de todas as regras de processo
	for _, analytics := range s.OrderProcessAnalytics {
		for _, employee := range analytics.EmployeePerformance {
			allEmployees = append(allEmployees, employee)
		}
	}

	// Ordena por score de eficiência (se implementado) ou por tempo médio
	// Aqui você pode implementar a lógica de ordenação desejada

	if limit > 0 && len(allEmployees) > limit {
		return allEmployees[:limit]
	}

	return allEmployees
}

func (s *Shift) IncrementCurrentOrder() {
	s.CurrentOrderNumber++
}

func (s *Shift) IsClosed() bool {
	return s.EndChange != nil
}

func (s *Shift) AddRedeem(redeem *Redeem) {
	s.Redeems = append(s.Redeems, *redeem)
}
