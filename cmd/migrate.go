package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

// MigrateCmd applies a raw SQL file to every tenant schema.
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Execute a SQL migration for every tenant schema",
	RunE: func(cmd *cobra.Command, _ []string) error {
		fileName, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		if fileName == "" {
			return fmt.Errorf("flag --file is required")
		}

		fullPath := fileName
		if !filepath.IsAbs(fileName) {
			fullPath = filepath.Join("bootstrap", "database", "migrations", fileName)
		}

		payload, err := os.ReadFile(fullPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("migration file %s not found (looked in %s)", fileName, fullPath)
			}
			return fmt.Errorf("failed to read migration file: %w", err)
		}

		sql := string(payload)
		if sql == "" {
			return fmt.Errorf("migration file %s is empty", fullPath)
		}

		cmd.Printf("connecting to database...\n")
		db := database.NewPostgreSQLConnection()

		ctx := context.Background()
		schemas, err := listTenantSchemas(ctx, db)
		if err != nil {
			return err
		}

		cmd.Printf("found %d tenant schemas\n", len(schemas))
		for _, schema := range schemas {
			cmd.Printf("applying migration to schema %s...\n", schema)
			if err := applyMigration(ctx, db, schema, sql); err != nil {
				return fmt.Errorf("schema %s: %w", schema, err)
			}
		}

		cmd.Println("migration applied to all schemas")
		return nil
	},
}

func listTenantSchemas(ctx context.Context, db *bun.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT schema_name
		FROM information_schema.schemata
		WHERE schema_name LIKE 'company_%'
		ORDER BY schema_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, err
		}
		schemas = append(schemas, schema)
	}

	return schemas, rows.Err()
}

func applyMigration(ctx context.Context, db *bun.DB, schemaName, sql string) error {
	schemaCtx := context.WithValue(ctx, model.Schema("schema"), schemaName)
	schemaCtx = database.WithTenantTransactionTimeout(schemaCtx, 0)
	schemaCtx, tx, cancel, err := database.GetTenantTransaction(schemaCtx, db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.ExecContext(schemaCtx, sql); err != nil {
		return err
	}

	return tx.Commit()
}
