package database

import (
	"fmt"

	"github.com/uptrace/bun"
	"golang.org/x/net/context"
)

const (
	addColumnSQL = `
		ALTER TABLE contacts ADD COLUMN IF NOT EXISTS ts tsvector
		GENERATED ALWAYS AS (
			setweight(to_tsvector('simple', coalesce(ddd::text, '') || coalesce(number::text, '')), 'A') ||
			setweight(to_tsvector('simple', coalesce(ddd::text, '')), 'B') ||
			setweight(to_tsvector('simple', coalesce(number::text, '')), 'B')
		) STORED;
	`
	createIndexSQL = "CREATE INDEX IF NOT EXISTS contacts_ts_idx ON contacts USING GIN(ts);"
)

func setupContactFtSearch(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, addColumnSQL); err != nil {
		return fmt.Errorf("failed to create tsvector column for contacts: %w", err)
	}

	if _, err := db.ExecContext(ctx, createIndexSQL); err != nil {
		return fmt.Errorf("failed to create index for contacts tsvector column: %w", err)
	}

	return nil
}
