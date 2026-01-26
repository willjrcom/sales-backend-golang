ALTER TABLE company_usage_costs
ADD COLUMN IF NOT EXISTS cycle_start TIMESTAMP,
ADD COLUMN IF NOT EXISTS cycle_end TIMESTAMP;

-- Index for date range queries (access control)
CREATE INDEX IF NOT EXISTS idx_company_usage_costs_cycle ON company_usage_costs(status, cycle_start, cycle_end);
