-- Migration: Add fiscal fields to companies table
-- Description: Adds fiscal invoice related fields to enable NFC-e emission
-- Date: 2026-01-06
-- Note: Transmitenota credentials are stored in ENV variables (TRANSMITENOTA_USUARIO, TRANSMITENOTA_SENHA)

ALTER TABLE companies 
ADD COLUMN IF NOT EXISTS fiscal_enabled BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS inscricao_estadual VARCHAR(20),
ADD COLUMN IF NOT EXISTS regime_tributario INTEGER DEFAULT 1,
ADD COLUMN IF NOT EXISTS cnae VARCHAR(10),
ADD COLUMN IF NOT EXISTS crt INTEGER DEFAULT 1;

COMMENT ON COLUMN companies.fiscal_enabled IS 'Feature toggle for fiscal invoice functionality';
COMMENT ON COLUMN companies.inscricao_estadual IS 'State registration (IE)';
COMMENT ON COLUMN companies.regime_tributario IS 'Tax regime: 1=Simples Nacional, 2=Simples Nacional Excesso, 3=Normal';
COMMENT ON COLUMN companies.cnae IS 'Economic activity code';
COMMENT ON COLUMN companies.crt IS 'Tax regime code for NFC-e (1, 2 or 3)';

