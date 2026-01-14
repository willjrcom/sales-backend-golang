-- Migration: Create fiscal_invoices table
-- Description: Table to store fiscal invoices (NFC-e, NF-e)
-- Date: 2026-01-06

CREATE TABLE IF NOT EXISTS fiscal_invoices (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    order_id UUID NOT NULL,
    chave_acesso VARCHAR(44) UNIQUE,
    numero INTEGER NOT NULL,
    serie INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'authorized', 'rejected', 'cancelled')),
    xml_path TEXT,
    pdf_path TEXT,
    protocolo VARCHAR(50),
    error_message TEXT,
    emitted_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    cancellation_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT unique_company_serie_numero UNIQUE (company_id, serie, numero)
);

CREATE INDEX IF NOT EXISTS idx_fiscal_invoices_company_id ON fiscal_invoices(company_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_fiscal_invoices_order_id ON fiscal_invoices(order_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_fiscal_invoices_chave_acesso ON fiscal_invoices(chave_acesso) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_fiscal_invoices_status ON fiscal_invoices(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_fiscal_invoices_created_at ON fiscal_invoices(created_at DESC) WHERE deleted_at IS NULL;

COMMENT ON TABLE fiscal_invoices IS 'Stores fiscal invoices (NFC-e, NF-e)';
COMMENT ON COLUMN fiscal_invoices.chave_acesso IS '44-character access key from SEFAZ';
COMMENT ON COLUMN fiscal_invoices.numero IS 'Invoice number (auto-incremented per company and series)';
COMMENT ON COLUMN fiscal_invoices.serie IS 'Invoice series (default: 1)';
COMMENT ON COLUMN fiscal_invoices.status IS 'Invoice status: pending, authorized, rejected, cancelled';
COMMENT ON COLUMN fiscal_invoices.xml_path IS 'Path to XML file';
COMMENT ON COLUMN fiscal_invoices.pdf_path IS 'Path to PDF file (DANFE)';
COMMENT ON COLUMN fiscal_invoices.protocolo IS 'Authorization protocol from SEFAZ';
