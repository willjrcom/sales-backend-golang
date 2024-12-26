package employeerepositorybun

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type EmployeeRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewEmployeeRepositoryBun(db *bun.DB) *EmployeeRepositoryBun {
	return &EmployeeRepositoryBun{db: db}
}

func (r *EmployeeRepositoryBun) CreateEmployee(ctx context.Context, c *employeeentity.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	// Create employee
	if _, err := r.db.NewInsert().Model(c).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepositoryBun) UpdateEmployee(ctx context.Context, p *employeeentity.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(p).Where("employee.id = ?", p.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepositoryBun) DeleteEmployee(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	// Delete employee
	if _, err := r.db.NewDelete().Model(&employeeentity.Employee{}).Where("employee.id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepositoryBun) GetEmployeeById(ctx context.Context, id string) (*employeeentity.Employee, error) {
	employee := &employeeentity.Employee{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(employee).Where("employee.id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return nil, err
	}

	user := &companyentity.User{}
	if err := r.db.NewSelect().Model(user).Where("u.id = ?", employee.UserID).Relation("Address").Relation("Contact").ExcludeColumn("hash").Scan(ctx); err != nil {
		return nil, err
	}

	employee.User = user
	return employee, nil
}

func (r *EmployeeRepositoryBun) GetEmployeeByUserID(ctx context.Context, userID string) (*employeeentity.Employee, error) {
	employee := &employeeentity.Employee{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(employee).Where("employee.user_id = ?", userID).Scan(ctx); err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *EmployeeRepositoryBun) GetAllEmployees(ctx context.Context) ([]employeeentity.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	employees := []employeeentity.Employee{}
	if err := r.db.NewSelect().Model(&employees).Scan(ctx); err != nil {
		return nil, err
	}

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return nil, err
	}
	// Extrair todos os UserIDs de uma vez
	userIDs := make([]uuid.UUID, len(employees))
	for i, employee := range employees {
		userIDs[i] = *employee.UserID
	}

	// Consultar todos os Users de uma vez
	users := []*companyentity.User{}
	if err := r.db.NewSelect().
		Model(&users).
		Where("u.id IN (?)", bun.In(userIDs)).
		Relation("Address").
		Relation("Contact").
		ExcludeColumn("hash").
		Scan(ctx); err != nil {
		return nil, err
	}
	// Mapear os usuários de volta para os funcionários
	userMap := make(map[uuid.UUID]*companyentity.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	for i := range employees {
		if user, exists := userMap[*employees[i].UserID]; exists {
			employees[i].User = user
		} else {
			// Tratar caso de usuário não encontrado
			employees[i].User = nil
		}
	}

	return employees, nil
}
