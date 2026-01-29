-- Add plan_type column to company_payments table
ALTER TABLE company_payments ADD COLUMN IF NOT EXISTS plan_type VARCHAR(20);

-- Add index for faster filtering by plan type
CREATE INDEX IF NOT EXISTS idx_company_payments_plan_type ON company_payments(plan_type);
