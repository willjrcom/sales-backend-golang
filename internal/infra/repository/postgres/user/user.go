package userrepositorybun

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
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

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	if _, err := tx.NewInsert().Model(user).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if user.Person.Contact != nil {
		if _, err := tx.NewDelete().Model(&personentity.Contact{}).Where("object_id = ?", user.ID).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

		// Create contact
		if _, err := tx.NewInsert().Model(user.Person.Contact).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if user.Person.Address != nil {
		if _, err := tx.NewDelete().Model(&addressentity.Address{}).Where("object_id = ?", user.ID).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

		// Create addresse
		if _, err := tx.NewInsert().Model(user.Person.Address).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
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

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	if _, err := tx.NewUpdate().Model(user).WherePK().Exec(context.Background()); err != nil {
		tx.Rollback()
		return err
	}

	if user.Person.Contact != nil {
		if _, err := tx.NewDelete().Model(&personentity.Contact{}).Where("object_id = ?", user.ID).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

		// Create contact
		if _, err := tx.NewInsert().Model(user.Person.Contact).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if user.Person.Address != nil {
		if _, err := tx.NewDelete().Model(&addressentity.Address{}).Where("object_id = ?", user.ID).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

		// Create addresse
		if _, err := tx.NewInsert().Model(user.Person.Address).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
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

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	if _, err := tx.NewSelect().Model(&companyentity.User{}).Where("u.email = ?", user.Email).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
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

func (r *UserRepositoryBun) ExistsUserByID(ctx context.Context, id uuid.UUID) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return false, err
	}

	user := &companyentity.User{}
	if err := r.db.NewSelect().Model(user).Where("u.id = ?", id).Column("id").Scan(ctx); err != nil {
		return false, err
	}

	return user.ID != uuid.Nil, nil
}
