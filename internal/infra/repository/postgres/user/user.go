package userrepositorybun

import (
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	userentity "github.com/willjrcom/sales-backend-go/internal/domain/user"
	"golang.org/x/net/context"
)

type UserRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewUserRepositoryBun(db *bun.DB) *UserRepositoryBun {
	return &UserRepositoryBun{db: db}
}

func (r *UserRepositoryBun) CreateUser(ctx context.Context, user *userentity.User) error {
	r.mu.Lock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
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

func (r *UserRepositoryBun) UpdateUser(ctx context.Context, user *userentity.User) error {
	r.mu.Lock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
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

func (r *UserRepositoryBun) DeleteUser(ctx context.Context, user *userentity.User) error {
	r.mu.Lock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		r.mu.Unlock()
		return err
	}

	_, err := r.db.NewDelete().Model(&userentity.User{}).Where("u.email = ? AND u.hash = ?", user.Email, user.Hash).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryBun) LoginUser(ctx context.Context, user *userentity.User) (*userentity.User, error) {
	r.mu.Lock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		r.mu.Unlock()
		return nil, err
	}

	err := r.db.NewSelect().
		Model(user).
		Where("u.email = ?", user.Email).
		Where("crypt(?, u.hash) = u.hash", user.Password).
		Limit(1).
		ExcludeColumn("hash").
		Scan(context.Background())

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryBun) GetIDByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), schemaentity.DEFAULT_SCHEMA)

	r.mu.Lock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		r.mu.Unlock()
		return uuid.Nil, err
	}

	user := &userentity.User{}
	err := r.db.NewSelect().Model(user).Where("u.email = ?", email).Column("id").Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}
