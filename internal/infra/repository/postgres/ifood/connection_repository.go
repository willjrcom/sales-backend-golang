package ifoodrepo

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ConnectionRepository struct {
	DB *bun.DB
}

func NewConnectionRepository(db *bun.DB) *ConnectionRepository {
	return &ConnectionRepository{DB: db}
}

// Upsert creates or updates an iFood connection by tenant schema.
func (r *ConnectionRepository) Upsert(ctx context.Context, conn *model.IfoodConnection) error {
	_, err := r.DB.NewInsert().Model(conn).
		On("CONFLICT (schema) DO UPDATE").
		Set("merchant_id = EXCLUDED.merchant_id").
		Set("client_id = EXCLUDED.client_id").
		Set("client_secret = EXCLUDED.client_secret").
		Set("access_token = EXCLUDED.access_token").
		Set("refresh_token = EXCLUDED.refresh_token").
		Set("token_expires_at = EXCLUDED.token_expires_at").
		Set("sandbox = EXCLUDED.sandbox").
		Exec(ctx)
	return err
}

func (r *ConnectionRepository) GetBySchema(ctx context.Context, schema string) (*model.IfoodConnection, error) {
	conn := &model.IfoodConnection{}
	if err := r.DB.NewSelect().Model(conn).Where("schema = ?", schema).Limit(1).Scan(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}
