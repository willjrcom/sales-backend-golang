-- Migration: Create company_usage_costs table
-- Description: Table to track monthly usage costs per company (subscription, NFC-e, etc.)
-- Date: 2026-01-06

CREATE TABLE IF NOT EXISTS company_usage_costs (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    cost_type VARCHAR(50) NOT NULL,
    description TEXT,
    amount DECIMAL(19, 4) NOT NULL,
    reference_id UUID,
    billing_month INTEGER NOT NULL CHECK (billing_month >= 1 AND billing_month <= 12),
    billing_year INTEGER NOT NULL CHECK (billing_year >= 2020),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_company_usage_costs_company_id ON company_usage_costs(company_id);
CREATE INDEX IF NOT EXISTS idx_company_usage_costs_billing_period ON company_usage_costs(company_id, billing_year, billing_month) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_company_usage_costs_type ON company_usage_costs(cost_type) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_company_usage_costs_reference ON company_usage_costs(reference_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE company_usage_costs IS 'Tracks monthly usage costs per company';
COMMENT ON COLUMN company_usage_costs.cost_type IS 'Type of cost: subscription, nfce, nfe, nfce_refund, nfe_refund';
COMMENT ON COLUMN company_usage_costs.amount IS 'Cost amount in BRL (uses decimal for precision)';
COMMENT ON COLUMN company_usage_costs.reference_id IS 'Optional reference to entity that generated the cost (e.g., fiscal_invoice.id)';
COMMENT ON COLUMN company_usage_costs.billing_month IS 'Month of billing (1-12)';
COMMENT ON COLUMN company_usage_costs.billing_year IS 'Year of billing';
