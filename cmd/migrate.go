package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

const migrationsDir = "bootstrap/database/migrations"
const publicMigrationsDir = "bootstrap/database/migrations/public"

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
			fullPath = filepath.Join(migrationsDir, fileName)
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

// MigrateCmd applies a raw SQL file to every tenant schema.
var PublicMigrateCmd = &cobra.Command{
	Use:   "public-migrate",
	Short: "Execute a SQL migration for public schema",
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
			fullPath = filepath.Join(publicMigrationsDir, fileName)
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

		schema := "public"
		cmd.Printf("applying migration to schema %s...\n", schema)
		if err := applyMigration(ctx, db, schema, sql); err != nil {
			return fmt.Errorf("schema %s: %w", schema, err)
		}

		cmd.Println("migration applied to schema public")
		return nil
	},
}

// MigrateAllCmd applies all pending SQL migrations to every tenant schema.
var MigrateAllCmd = &cobra.Command{
	Use:   "migrate-all",
	Short: "Execute all pending SQL migrations for every tenant schema",
	RunE: func(cmd *cobra.Command, _ []string) error {
		cmd.Printf("connecting to database...\n")
		db := database.NewPostgreSQLConnection()
		ctx := context.Background()

		// Cria tabela de controle de migrações no schema public
		if err := ensureMigrationsTable(ctx, db); err != nil {
			return fmt.Errorf("failed to create migrations table: %w", err)
		}

		// Lista todos os schemas de tenant
		schemas, err := listTenantSchemas(ctx, db)
		if err != nil {
			return err
		}

		if len(schemas) == 0 {
			cmd.Println("no tenant schemas found, skipping migrations")
			return nil
		}

		cmd.Printf("found %d tenant schemas\n", len(schemas))

		// Lista todas as migrações disponíveis
		migrations, err := listMigrationFiles()
		if err != nil {
			return fmt.Errorf("failed to list migration files: %w", err)
		}

		if len(migrations) == 0 {
			cmd.Println("no migration files found")
			return nil
		}

		cmd.Printf("found %d migration files\n", len(migrations))

		// Para cada schema, aplica migrações pendentes
		for _, schema := range schemas {
			applied, err := getAppliedMigrations(ctx, db, schema)
			if err != nil {
				return fmt.Errorf("schema %s: failed to get applied migrations: %w", schema, err)
			}

			pending := filterPendingMigrations(migrations, applied)
			if len(pending) == 0 {
				cmd.Printf("schema %s: all migrations already applied\n", schema)
				continue
			}

			cmd.Printf("schema %s: applying %d pending migrations...\n", schema, len(pending))

			for _, migration := range pending {
				cmd.Printf("  -> %s\n", migration)

				fullPath := filepath.Join(migrationsDir, migration)
				sql, err := readMigrationFile(fullPath)
				if err != nil {
					return fmt.Errorf("schema %s: %w", schema, err)
				}

				if err := applyMigration(ctx, db, schema, sql); err != nil {
					return fmt.Errorf("schema %s: migration %s failed: %w", schema, migration, err)
				}

				if err := recordMigration(ctx, db, schema, migration); err != nil {
					return fmt.Errorf("schema %s: failed to record migration %s: %w", schema, migration, err)
				}
			}
		}

		cmd.Println("all migrations applied successfully")
		return nil
	},
}

func ensureMigrationsTable(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS public.schema_migrations (
			id SERIAL PRIMARY KEY,
			schema_name VARCHAR(255) NOT NULL,
			migration_name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(schema_name, migration_name)
		)
	`)
	return err
}

func getAppliedMigrations(ctx context.Context, db *bun.DB, schemaName string) (map[string]bool, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT migration_name FROM public.schema_migrations WHERE schema_name = ?
	`, schemaName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied[name] = true
	}

	return applied, rows.Err()
}

func recordMigration(ctx context.Context, db *bun.DB, schemaName, migrationName string) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO public.schema_migrations (schema_name, migration_name, applied_at)
		VALUES (?, ?, ?)
		ON CONFLICT (schema_name, migration_name) DO NOTHING
	`, schemaName, migrationName, time.Now().UTC())
	return err
}

func listMigrationFiles() ([]string, error) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var migrations []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".sql") {
			migrations = append(migrations, entry.Name())
		}
	}

	// Ordena por nome (que deve incluir timestamp no prefixo)
	sort.Strings(migrations)
	return migrations, nil
}

func listPublicMigrationFiles() ([]string, error) {
	entries, err := os.ReadDir(publicMigrationsDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var migrations []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".sql") {
			migrations = append(migrations, entry.Name())
		}
	}

	// Ordena por nome (que deve incluir timestamp no prefixo)
	sort.Strings(migrations)
	return migrations, nil
}

func filterPendingMigrations(all []string, applied map[string]bool) []string {
	var pending []string
	for _, m := range all {
		if !applied[m] {
			pending = append(pending, m)
		}
	}
	return pending
}

func readMigrationFile(fullPath string) (string, error) {
	payload, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read migration file %s: %w", fullPath, err)
	}

	sql := string(payload)
	if sql == "" {
		return "", fmt.Errorf("migration file %s is empty", fullPath)
	}

	return sql, nil
}

func listTenantSchemas(ctx context.Context, db *bun.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT schema_name
		FROM information_schema.schemata
		WHERE schema_name LIKE 'company_%' OR schema_name LIKE 'loja_%'
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

// MigrateAllCmd applies all pending SQL migrations to every tenant schema.
var PublicMigrateAllCmd = &cobra.Command{
	Use:   "public-migrate-all",
	Short: "Execute all pending SQL migrations for public schema",
	RunE: func(cmd *cobra.Command, _ []string) error {
		cmd.Printf("connecting to database...\n")
		db := database.NewPostgreSQLConnection()
		ctx := context.Background()

		// Cria tabela de controle de migrações no schema public
		if err := ensureMigrationsTable(ctx, db); err != nil {
			return fmt.Errorf("failed to create migrations table: %w", err)
		}

		// Lista todas as migrações disponíveis
		migrations, err := listPublicMigrationFiles()
		if err != nil {
			return fmt.Errorf("failed to list migration files: %w", err)
		}

		if len(migrations) == 0 {
			cmd.Println("no migration files found")
			return nil
		}

		cmd.Printf("found %d migration files\n", len(migrations))

		// Para cada schema, aplica migrações pendentes
		applied, err := getAppliedMigrations(ctx, db, "public")
		if err != nil {
			return fmt.Errorf("schema public: failed to get applied migrations: %w", err)
		}

		pending := filterPendingMigrations(migrations, applied)
		if len(pending) == 0 {
			cmd.Printf("schema public: all migrations already applied")
			return nil
		}

		cmd.Printf("schema public: applying %d pending migrations...\n", len(pending))

		for _, migration := range pending {
			cmd.Printf("  -> %s\n", migration)

			fullPath := filepath.Join(publicMigrationsDir, migration)
			sql, err := readMigrationFile(fullPath)
			if err != nil {
				return fmt.Errorf("schema public: %w", err)
			}

			if err := applyMigration(ctx, db, "public", sql); err != nil {
				return fmt.Errorf("schema public: migration %s failed: %w", migration, err)
			}

			if err := recordMigration(ctx, db, "public", migration); err != nil {
				return fmt.Errorf("schema public: failed to record migration %s: %w", migration, err)
			}
		}

		cmd.Println("all migrations applied successfully")
		return nil
	},
}
