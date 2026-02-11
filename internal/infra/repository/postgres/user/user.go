package userrepositorybun

import (
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	"golang.org/x/net/context"
)

type UserRepositoryBun struct {
	db *bun.DB
}

func NewUserRepositoryBun(db *bun.DB) model.UserRepository {
	return &UserRepositoryBun{db: db}
}

func (r *UserRepositoryBun) CreateUser(ctx context.Context, user *model.User) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(user).Exec(ctx); err != nil {
		return err
	}

	if user.PublicPerson.Contact != nil {
		res, err := tx.NewUpdate().Model(user.PublicPerson.Contact).Where("id = ?", user.PublicPerson.Contact.ID).Exec(ctx)
		if err != nil {
			return err
		}

		if rows, _ := res.RowsAffected(); rows == 0 {
			if _, err := tx.NewInsert().Model(user.PublicPerson.Contact).Exec(ctx); err != nil {
				return err
			}
		}
	}

	if user.PublicPerson.Address != nil {
		res, err := tx.NewUpdate().Model(user.PublicPerson.Address).Where("id = ?", user.PublicPerson.Address.ID).Exec(ctx)
		if err != nil {
			return err
		}

		if rows, _ := res.RowsAffected(); rows == 0 {
			if _, err := tx.NewInsert().Model(user.PublicPerson.Address).Exec(ctx); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryBun) GetUser(ctx context.Context, email string) (*model.User, error) {

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

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

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryBun) UpdateUserPassword(ctx context.Context, user *model.User) error {

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().Model(user).WherePK().Column("hash").Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryBun) UpdateUser(ctx context.Context, user *model.User) error {

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().Model(user).WherePK().ExcludeColumn("hash").Exec(ctx); err != nil {
		return err
	}

	if user.Contact != nil {
		res, err := tx.NewUpdate().Model(user.Contact).WherePK().Exec(ctx)
		if err != nil {
			return err
		}

		if rows, _ := res.RowsAffected(); rows == 0 {
			if _, err := tx.NewInsert().Model(user.Contact).Exec(ctx); err != nil {
				return err
			}
		}
	}

	if user.Address != nil {
		res, err := tx.NewUpdate().Model(user.Address).WherePK().Exec(ctx)
		if err != nil {
			return err
		}

		if rows, _ := res.RowsAffected(); rows == 0 {
			if _, err := tx.NewInsert().Model(user.Address).Exec(ctx); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
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

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewDelete().Model(userLogged).Where("id = ?", userLogged.ID).Exec(ctx); err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model((&model.Address{})).Where("object_id = ?", userLogged.ID).Exec(ctx); err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model((&model.Contact{})).Where("object_id = ?", userLogged.ID).Exec(ctx); err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(userLogged).Where("id = ?", userLogged.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// ListPublicUsers returns every user stored in the public schema with limited fields.
func (r *UserRepositoryBun) ListPublicUsers(ctx context.Context) ([]model.User, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	users := []model.User{}
	if err := tx.NewSelect().
		Model(&users).
		Column("id", "name", "email", "cpf").Relation("Companies").
		Order("name ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepositoryBun) LoginUser(ctx context.Context, user *model.User) (*model.User, error) {

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

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

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryBun) GetIDByEmail(ctx context.Context, email string) (*uuid.UUID, error) {

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	user := &model.User{}
	if err := tx.NewSelect().Model(user).Where("u.email = ?", email).Column("id").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &user.ID, nil
}

func (r *UserRepositoryBun) GetIDByEmailOrCPF(ctx context.Context, email string, cpf string) (*uuid.UUID, string, error) {

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, "", err
	}

	defer cancel()
	defer tx.Rollback()

	user := &model.User{}
	if err := tx.NewSelect().Model(user).Where("u.email = ? or u.cpf = ?", email, cpf).Column("id", "email", "cpf").Scan(ctx); err != nil {
		return nil, "", err
	}

	if err := tx.Commit(); err != nil {
		return nil, "", err
	}

	key := ""
	if user.Email == email {
		key = "email"
	}

	if user.Cpf == cpf {
		key = "cpf"
	}

	return &user.ID, key, nil
}

func (r *UserRepositoryBun) GetUserByID(ctx context.Context, id uuid.UUID, withCompanies bool) (*model.User, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	user := &model.User{}
	query := tx.NewSelect().Model(user).Where("u.id = ?", id).Relation("Address").Relation("Contact")

	if withCompanies {
		query = query.Relation("Companies")
	}

	if err := query.ExcludeColumn("hash").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryBun) GetByCPF(ctx context.Context, cpf string) (*model.User, error) {

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	user := &model.User{}
	if err := tx.NewSelect().Model(user).Where("u.cpf = ?", cpf).Relation("Address").Relation("Contact").ExcludeColumn("hash").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryBun) ExistsUserByID(ctx context.Context, id uuid.UUID) (bool, error) {

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return false, err
	}

	defer cancel()
	defer tx.Rollback()

	user := &model.User{}
	if err := tx.NewSelect().Model(user).Where("u.id = ?", id).Column("id").Scan(ctx); err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}
	return user.ID != uuid.Nil, nil
}
