package database

import (
	"errors"

	"github.com/uptrace/bun"
	"golang.org/x/net/context"
)

func SetupContactFtSearch(ctx context.Context, db *bun.DB) error {
	if err := ChangeSchema(ctx, db); err != nil {
		return err
	}

	column := `
		ALTER TABLE contacts ADD COLUMN IF NOT EXISTS ts tsvector
		GENERATED ALWAYS AS (
			setweight(to_tsvector('simple', coalesce(ddd::text, '') || coalesce(number::text, '')), 'A') ||
			setweight(to_tsvector('simple', coalesce(ddd::text, '')), 'B') ||
			setweight(to_tsvector('simple', coalesce(number::text, '')), 'B')
		) STORED;
	`

	if _, err := db.ExecContext(ctx, column); err != nil {
		return errors.New("Failed to create tsvector column for contacts" + err.Error())
	}

	index := "CREATE INDEX IF NOT EXISTS contacts_ts_idx ON contacts USING GIN(ts);"

	if _, err := db.ExecContext(ctx, index); err != nil {
		return errors.New("Failed to create index for contacts tsvector column" + err.Error())
	}

	return nil
}
