package processentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Process struct {
	entity.Entity
	bun.BaseModel `bun:"table:processes"`
	ProcessTimeLogs
	ProcessCommonAttributes
}

type ProcessCommonAttributes struct {
	EmployeeID uuid.UUID                `bun:"column:employee_id,type:uuid,notnull" json:"employee_id"`
	Employee   *employeeentity.Employee `bun:"rel:belongs-to" json:"employee,omitempty"`
	ItemID     uuid.UUID                `bun:"item_id,type:uuid,notnull" json:"item_id"`
	ProcessID  uuid.UUID                `bun:"process_id,type:uuid,notnull" json:"process_id"`
}

type ProcessTimeLogs struct {
	StartedAt  time.Time  `bun:"started_at" json:"started_at,omitempty"`
	FinishedAt *time.Time `bun:"finished_at" json:"finished_at,omitempty"`
}

func NewProcess(processCommonAttributes ProcessCommonAttributes) *Process {
	return &Process{
		Entity:                  entity.NewEntity(),
		ProcessCommonAttributes: processCommonAttributes,
		ProcessTimeLogs: ProcessTimeLogs{
			StartedAt: time.Now(),
		},
	}
}

func (p *Process) FinishProcess() {
	p.FinishedAt = &time.Time{}
	*p.FinishedAt = time.Now()
}
