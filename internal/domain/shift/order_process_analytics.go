package shiftentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
)

// OrderProcessAnalytics representa as métricas agregadas de processos de produção
type OrderProcessAnalytics struct {
	ProcessRuleID       uuid.UUID
	ProcessRuleName     string
	TotalProcesses      int
	CompletedProcesses  int
	CanceledProcesses   int
	AverageProcessTime  time.Duration
	TotalProcessTime    time.Duration
	TotalPausedTime     time.Duration
	TotalPausedCount    int
	AverageQueueTime    time.Duration
	TotalQueueTime      time.Duration
	ProcessesByStatus   map[orderprocessentity.StatusProcess]int
	CategoriesProcessed map[uuid.UUID]CategoryProcessMetrics
	EmployeePerformance map[uuid.UUID]EmployeeProcessMetrics
}

// CategoryProcessMetrics representa métricas por categoria de produto
type CategoryProcessMetrics struct {
	CategoryID         uuid.UUID
	CategoryName       string
	TotalProcessed     int
	AverageProcessTime time.Duration
	TotalProcessTime   time.Duration
	CanceledCount      int
}

// EmployeeProcessMetrics representa métricas por funcionário
type EmployeeProcessMetrics struct {
	EmployeeID         uuid.UUID
	EmployeeName       string
	TotalProcessed     int
	AverageProcessTime time.Duration
	TotalProcessTime   time.Duration
	EfficiencyScore    decimal.Decimal // Tempo médio vs tempo esperado
}

// QueueAnalytics representa métricas agregadas das filas
type QueueAnalytics struct {
	ProcessRuleID      uuid.UUID
	ProcessRuleName    string
	TotalQueued        int
	AverageQueueTime   time.Duration
	MaxQueueTime       time.Duration
	MinQueueTime       time.Duration
	TotalQueueTime     time.Duration
	PeakQueueLength    int
	AverageQueueLength int
}

// NewOrderProcessAnalytics cria uma nova instância de analytics
func NewOrderProcessAnalytics(processRuleID uuid.UUID, processRuleName string) *OrderProcessAnalytics {
	return &OrderProcessAnalytics{
		ProcessRuleID:       processRuleID,
		ProcessRuleName:     processRuleName,
		ProcessesByStatus:   make(map[orderprocessentity.StatusProcess]int),
		CategoriesProcessed: make(map[uuid.UUID]CategoryProcessMetrics),
		EmployeePerformance: make(map[uuid.UUID]EmployeeProcessMetrics),
	}
}

// AddProcess adiciona um processo às métricas
func (a *OrderProcessAnalytics) AddProcess(process *orderprocessentity.OrderProcess, employeeName string) {
	a.TotalProcesses++
	a.ProcessesByStatus[process.Status]++

	// Conta processos por status
	switch process.Status {
	case orderprocessentity.ProcessStatusFinished:
		a.CompletedProcesses++
	case orderprocessentity.ProcessStatusCanceled:
		a.CanceledProcesses++
	}

	// Calcula tempos
	if process.FinishedAt != nil && process.StartedAt != nil {
		processTime := process.Duration
		a.TotalProcessTime += processTime
		if a.CompletedProcesses > 0 {
			a.AverageProcessTime = time.Duration(int64(a.TotalProcessTime) / int64(a.CompletedProcesses))
		}
	}

	// Tempo de pausa
	if process.TotalPaused > 0 {
		a.TotalPausedCount += int(process.TotalPaused)
		// Aqui você pode calcular o tempo total de pausa se tiver os dados
	}

	// Métricas por categoria (se disponível)
	if process.GroupItem != nil && process.GroupItem.Category != nil {
		categoryID := process.GroupItem.CategoryID
		categoryMetrics, exists := a.CategoriesProcessed[categoryID]

		if !exists {
			categoryMetrics = CategoryProcessMetrics{
				CategoryID:   categoryID,
				CategoryName: process.GroupItem.Category.Name,
			}
		}

		categoryMetrics.TotalProcessed++
		if process.Status == orderprocessentity.ProcessStatusCanceled {
			categoryMetrics.CanceledCount++
		}

		if process.FinishedAt != nil && process.StartedAt != nil {
			categoryMetrics.TotalProcessTime += process.Duration
			if categoryMetrics.TotalProcessed > 0 {
				categoryMetrics.AverageProcessTime = time.Duration(int64(categoryMetrics.TotalProcessTime) / int64(categoryMetrics.TotalProcessed))
			}
		}

		a.CategoriesProcessed[categoryID] = categoryMetrics
	}

	// Métricas por funcionário
	if process.EmployeeID != nil {
		employeeID := *process.EmployeeID
		employeeMetrics, exists := a.EmployeePerformance[employeeID]

		if !exists {
			employeeMetrics = EmployeeProcessMetrics{
				EmployeeID:   employeeID,
				EmployeeName: employeeName,
			}
		}

		employeeMetrics.TotalProcessed++
		if process.FinishedAt != nil && process.StartedAt != nil {
			employeeMetrics.TotalProcessTime += process.Duration
			if employeeMetrics.TotalProcessed > 0 {
				employeeMetrics.AverageProcessTime = time.Duration(int64(employeeMetrics.TotalProcessTime) / int64(employeeMetrics.TotalProcessed))
			}
		}

		a.EmployeePerformance[employeeID] = employeeMetrics
	}
}

// AddQueue adiciona uma fila às métricas
func (a *OrderProcessAnalytics) AddQueue(queue *orderprocessentity.OrderQueue) {
	if queue.LeftAt != nil {
		queueTime := queue.Duration
		a.TotalQueueTime += queueTime
		if a.CompletedProcesses > 0 {
			a.AverageQueueTime = time.Duration(int64(a.TotalQueueTime) / int64(a.CompletedProcesses))
		}
	}
}

// GetEfficiencyScore calcula o score de eficiência baseado no tempo médio vs tempo esperado
func (a *OrderProcessAnalytics) GetEfficiencyScore(expectedAverageTime time.Duration) decimal.Decimal {
	if a.AverageProcessTime == 0 || expectedAverageTime == 0 {
		return decimal.Zero
	}

	// Score baseado na relação entre tempo esperado e tempo real
	// Quanto menor o tempo real vs esperado, maior o score
	ratio := float64(expectedAverageTime) / float64(a.AverageProcessTime)
	return decimal.NewFromFloat(ratio * 100) // Score de 0-100
}

// NewQueueAnalytics cria uma nova instância de analytics de fila
func NewQueueAnalytics(processRuleID uuid.UUID, processRuleName string) *QueueAnalytics {
	return &QueueAnalytics{
		ProcessRuleID:   processRuleID,
		ProcessRuleName: processRuleName,
	}
}

// AddQueue adiciona uma fila às métricas
func (q *QueueAnalytics) AddQueue(queue *orderprocessentity.OrderQueue) {
	q.TotalQueued++

	if queue.LeftAt != nil {
		queueTime := queue.Duration
		q.TotalQueueTime += queueTime
		if q.TotalQueued > 0 {
			q.AverageQueueTime = time.Duration(int64(q.TotalQueueTime) / int64(q.TotalQueued))
		}

		// Atualiza min/max
		if q.MaxQueueTime == 0 || queueTime > q.MaxQueueTime {
			q.MaxQueueTime = queueTime
		}
		if q.MinQueueTime == 0 || queueTime < q.MinQueueTime {
			q.MinQueueTime = queueTime
		}
	}
}
