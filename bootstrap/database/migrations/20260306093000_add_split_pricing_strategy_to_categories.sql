ALTER TABLE product_categories
    ADD COLUMN IF NOT EXISTS split_pricing_strategy TEXT NOT NULL DEFAULT 'highest_item';
