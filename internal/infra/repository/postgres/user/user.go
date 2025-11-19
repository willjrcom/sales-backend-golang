package userrepositorybun

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	"golang.org/x/net/context"
)

type UserRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewUserRepositoryBun(db *bun.DB) model.UserRepository {
	return &UserRepositoryBun{db: db}
}

func (r *UserRepositoryBun) CreateUser(ctx context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewInsert().Model(user).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if user.PublicPerson.Contact != nil {
		if _, err := tx.NewDelete().Model(&model.Contact{}).Where("object_id = ?", user.ID).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

		// Create contact
		if _, err := tx.NewInsert().Model(user.PublicPerson.Contact).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if user.PublicPerson.Address != nil {
		if _, err := tx.NewDelete().Model(&model.Address{}).Where("object_id = ?", user.ID).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

		// Create addresse
		if _, err := tx.NewInsert().Model(user.PublicPerson.Address).Exec(ctx); err != nil {
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

func (r *UserRepositoryBun) GetUser(ctx context.Context, email string) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	user := &model.User{}
	if err := tx.NewSelect().
		Model(user).
		Where("u.email = ?", email).
		Relation("Companies").Relation("Contact").Relation("Address").
		Limit(1).
		ExcludeColumn("hash").
		Scan(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryBun) UpdateUserPassword(ctx context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewUpdate().Model(user).WherePK().Column("hash").Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryBun) UpdateUser(ctx context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewUpdate().Model(user).WherePK().ExcludeColumn("hash").Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if user.Contact != nil {
		if _, err := tx.NewUpdate().Model(user.Contact).WherePK().Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if user.Address != nil {
		if _, err := tx.NewUpdate().Model(user.Address).WherePK().Exec(ctx); err != nil {
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

func (r *UserRepositoryBun) LoginAndDeleteUser(ctx context.Context, user *model.User) error {
	userLogged, err := r.LoginUser(ctx, user)
	if err != nil {
		return err
	}

	if userLogged == nil {
		return errors.New("invalid email or password")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	userToDelete := &model.User{}

	if _, err := tx.NewSelect().Model(userToDelete).Where("u.email = ?", user.Email).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model((&model.Person{})).Where("object_id = ?", userToDelete.ID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model((&model.Address{})).Where("object_id = ?", userToDelete.ID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model((&model.Contact{})).Where("object_id = ?", userToDelete.ID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(userToDelete).Where("id = ?", userToDelete.ID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryBun) LoginUser(ctx context.Context, user *model.User) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().
		Model(user).
		Where("u.email = ?", user.Email).
		Where("crypt(?, u.hash) = u.hash", user.Password).
		Relation("Companies").Relation("Contact").Relation("Address").
		Limit(1).
		ExcludeColumn("hash").
		Scan(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryBun) GetIDByEmail(ctx context.Context, email string) (*uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	user := &model.User{}
	if err := tx.NewSelect().Model(user).Where("u.email = ?", email).Column("id").Scan(ctx); err != nil {
		return nil, err
	}

	return &user.ID, nil
}

func (r *UserRepositoryBun) GetIDByEmailOrCPF(ctx context.Context, email string, cpf string) (*uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	user := &model.User{}
	if err := tx.NewSelect().Model(user).Where("u.email = ? or u.cpf = ?", email, cpf).Column("id").Scan(ctx); err != nil {
		return nil, err
	}

	return &user.ID, nil
}

func (r *UserRepositoryBun) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	user := &model.User{}
	if err := tx.NewSelect().Model(user).Where("u.id = ?", id).Relation("Companies").Relation("Address").Relation("Contact").ExcludeColumn("hash").Scan(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryBun) GetByCPF(ctx context.Context, cpf string) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	user := &model.User{}
	if err := tx.NewSelect().Model(user).Where("u.cpf = ?", cpf).Relation("Address").Relation("Contact").ExcludeColumn("hash").Scan(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryBun) ExistsUserByID(ctx context.Context, id uuid.UUID) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return false, err
	}

	user := &model.User{}
	if err := tx.NewSelect().Model(user).Where("u.id = ?", id).Column("id").Scan(ctx); err != nil {
		return false, err
	}

	return user.ID != uuid.Nil, nil
}
