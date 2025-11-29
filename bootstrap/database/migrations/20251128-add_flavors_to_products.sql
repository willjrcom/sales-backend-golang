ALTER TABLE products
    ADD COLUMN IF NOT EXISTS flavors jsonb NOT NULL DEFAULT '[]'::jsonb;
