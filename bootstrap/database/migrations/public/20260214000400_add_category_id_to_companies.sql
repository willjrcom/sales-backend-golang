-- Create company_categories table
CREATE TABLE IF NOT EXISTS company_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    image_path TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Add category_id column to companies table
ALTER TABLE companies ADD COLUMN IF NOT EXISTS category_id UUID;

-- Add foreign key constraint to company_categories table (with idempotent check)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_companies_category'
    ) THEN
        ALTER TABLE companies ADD CONSTRAINT fk_companies_category 
            FOREIGN KEY (category_id) REFERENCES company_categories(id) ON DELETE SET NULL;
    END IF;
END $$;

-- Add index for better query performance
CREATE INDEX IF NOT EXISTS idx_companies_category_id ON companies(category_id);
