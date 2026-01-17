-- Add is_active column to process_rules table
-- This migration adds the missing is_active column

ALTER TABLE process_rules ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT true;
