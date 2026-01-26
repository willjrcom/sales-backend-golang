ALTER TABLE company_usage_costs
ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
ADD COLUMN IF NOT EXISTS payment_id UUID,
ADD COLUMN IF NOT EXISTS original_amount DECIMAL(19,4);

-- Rename columns to simplify
ALTER TABLE company_usage_costs
RENAME COLUMN billing_month TO month;

ALTER TABLE company_usage_costs
RENAME COLUMN billing_year TO year;

-- Add description to payments
ALTER TABLE company_payments
ADD COLUMN IF NOT EXISTS description TEXT;

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_company_usage_costs_payment_id ON company_usage_costs(payment_id);
CREATE INDEX IF NOT EXISTS idx_company_usage_costs_status_month_year ON company_usage_costs(status, month, year);
