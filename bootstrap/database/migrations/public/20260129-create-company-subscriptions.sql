-- Create company_subscriptions table
CREATE TABLE IF NOT EXISTS company_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    payment_id UUID REFERENCES company_payments(id), -- Nullable for manual grants/trials
    plan_type TEXT NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add current_plan to companies (snapshot/cache)
ALTER TABLE companies ADD COLUMN IF NOT EXISTS current_plan TEXT DEFAULT 'free';
