ALTER TABLE company_usage_costs
ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
ADD COLUMN IF NOT EXISTS payment_id UUID,
ADD COLUMN IF NOT EXISTS original_amount DECIMAL(19,4);

-- Add description to payments
ALTER TABLE company_payments
ADD COLUMN IF NOT EXISTS description TEXT;
