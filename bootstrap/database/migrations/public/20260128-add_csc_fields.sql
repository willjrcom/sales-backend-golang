-- Add CSC fields to fiscal_settings
ALTER TABLE fiscal_settings ADD COLUMN IF NOT EXISTS csc_production_id TEXT;
ALTER TABLE fiscal_settings ADD COLUMN IF NOT EXISTS csc_production_code TEXT;
ALTER TABLE fiscal_settings ADD COLUMN IF NOT EXISTS csc_homologation_id TEXT;
ALTER TABLE fiscal_settings ADD COLUMN IF NOT EXISTS csc_homologation_code TEXT;
