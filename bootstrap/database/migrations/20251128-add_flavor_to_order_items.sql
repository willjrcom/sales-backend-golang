ALTER TABLE order_items
    ADD COLUMN IF NOT EXISTS flavor text;
