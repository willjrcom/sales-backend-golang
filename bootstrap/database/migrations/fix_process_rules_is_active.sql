-- Add is_active column to process_rules in all tenant schemas
-- Run this manually with: go run main.go exec-sql fix_process_rules_is_active.sql

DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT nspname 
        FROM pg_namespace 
        WHERE nspname LIKE 'company_%'
    LOOP
        EXECUTE format('
            ALTER TABLE %I.process_rules 
            ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT true;
        ', schema_name);
        
        RAISE NOTICE 'Added is_active to %.process_rules', schema_name;
    END LOOP;
END $$;
