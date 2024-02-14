package userrepositorybun

import (
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"golang.org/x/net/context"
)

type UserRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewUserRepositoryBun(db *bun.DB) *UserRepositoryBun {
	return &UserRepositoryBun{db: db}
}

func (r *UserRepositoryBun) CreateUser(ctx context.Context, user *companyentity.User) error {
	r.mu.Lock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		r.mu.Unlock()
		return err
	}

	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryBun) UpdateUser(ctx context.Context, user *companyentity.User) error {
	r.mu.Lock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		r.mu.Unlock()
		return err
	}

	_, err := r.db.NewUpdate().Model(user).WherePK().Exec(context.Background())
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryBun) DeleteUser(ctx context.Context, user *companyentity.User) error {
	r.mu.Lock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		r.mu.Unlock()
		return err
	}

	_, err := r.db.NewDelete().Model(&companyentity.User{}).Where("u.email = ? AND u.hash = ?", user.Email, user.Hash).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryBun) LoginUser(ctx context.Context, user *companyentity.User) (*companyentity.User, error) {
	r.mu.Lock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		r.mu.Unlock()
		return nil, err
	}

	err := r.db.NewSelect().
		Model(user).
		Where("u.email = ?", user.Email).
		Where("crypt(?, u.hash) = u.hash", user.Password).
		Relation("CompanyToUsers").
		Limit(1).
		ExcludeColumn("hash").
		Scan(context.Background())

	if err != nil {
		r.mu.Unlock()
		return nil, err
	}

	for _, ctu := range user.CompanyToUsers {
		company := &companyentity.CompanyWithUsers{}
		if err = r.db.NewSelect().Model(company).Where("id = ?", ctu.CompanyWithUsersID).Scan(ctx); err != nil {
			r.mu.Unlock()
			return nil, err
		}

		company.Address = nil
		user.Companies = append(user.Companies, *company)
	}

	r.mu.Unlock()
	return user, nil
}

func (r *UserRepositoryBun) GetIDByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	r.mu.Lock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		r.mu.Unlock()
		return uuid.Nil, err
	}

	user := &companyentity.User{}
	err := r.db.NewSelect().Model(user).Where("u.email = ?", email).Column("id").Scan(ctx)

	r.mu.Unlock()
	return user.ID, err
}
