package database

import (
	"fmt"

	"github.com/uptrace/bun"
	"golang.org/x/net/context"
)

const (
	addReadyAtOrdersSQL         = "ALTER TABLE orders ADD COLUMN IF NOT EXISTS ready_at timestamptz;"
	addSubscriptionExpiresAtSQL = "ALTER TABLE companies ADD COLUMN IF NOT EXISTS subscription_expires_at timestamptz;"
)

// setupPrivateMigrations ensures every table that expects a ready_at column has it.
func setupPrivateMigrations(ctx context.Context, tx *bun.Tx) error {
	if _, err := tx.ExecContext(ctx, addReadyAtOrdersSQL); err != nil {
		return fmt.Errorf("failed to add ready_at to orders: %w", err)
	}
	if _, err := tx.ExecContext(ctx, addSubscriptionExpiresAtSQL); err != nil {
		return fmt.Errorf("failed to add subscription_expires_at to companies: %w", err)
	}

	return nil
}

func setupPublicMigrations(ctx context.Context, tx *bun.Tx) error {
	if _, err := tx.ExecContext(ctx, addSubscriptionExpiresAtSQL); err != nil {
		return fmt.Errorf("failed to add subscription_expires_at to public companies: %w", err)
	}
	return nil
}
