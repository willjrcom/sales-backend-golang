package userrepositorybun

import (
	"errors"
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
	defer r.mu.Unlock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewInsert().Model(user).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryBun) UpdateUser(ctx context.Context, user *companyentity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(user).WherePK().Exec(context.Background()); err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryBun) LoginAndDeleteUser(ctx context.Context, user *companyentity.User) error {
	userLogged, err := r.LoginUser(ctx, user)
	if err != nil {
		return err
	}

	if userLogged == nil {
		return errors.New("invalid email or password")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&companyentity.User{}).Where("u.email = ?", user.Email).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryBun) LoginUser(ctx context.Context, user *companyentity.User) (*companyentity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().
		Model(user).
		Where("u.email = ?", user.Email).
		Where("crypt(?, u.hash) = u.hash", user.Password).
		Relation("CompanyToUsers").
		Limit(1).
		ExcludeColumn("hash").
		Scan(context.Background()); err != nil {
		return nil, err
	}

	for _, ctu := range user.CompanyToUsers {
		company := &companyentity.CompanyWithUsers{}
		if err := r.db.NewSelect().Model(company).Where("id = ?", ctu.CompanyWithUsersID).Scan(ctx); err != nil {
			return nil, err
		}

		company.Address = nil
		user.Companies = append(user.Companies, *company)
	}

	return user, nil
}

func (r *UserRepositoryBun) GetIDByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return uuid.Nil, err
	}

	user := &companyentity.User{}
	err := r.db.NewSelect().Model(user).Where("u.email = ?", email).Column("id").Scan(ctx)

	return user.ID, err
}
