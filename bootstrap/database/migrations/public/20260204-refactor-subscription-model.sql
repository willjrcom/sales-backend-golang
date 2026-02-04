-- Add new columns for subscription refactor
ALTER TABLE company_subscriptions
ADD COLUMN IF NOT EXISTS preapproval_id TEXT,
ADD COLUMN IF NOT EXISTS external_reference TEXT,
ADD COLUMN IF NOT EXISTS is_canceled BOOLEAN DEFAULT FALSE;

-- Ensure preapproval_id is unique to prevent duplicate subscriptions for same contract
CREATE UNIQUE INDEX IF NOT EXISTS idx_company_subscriptions_preapproval_id ON company_subscriptions (preapproval_id) WHERE preapproval_id IS NOT NULL;
