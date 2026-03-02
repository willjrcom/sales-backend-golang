-- Rename columns in order_items
ALTER TABLE order_items RENAME COLUMN IF EXISTS price TO sub_total;
ALTER TABLE order_items RENAME COLUMN IF EXISTS total_price TO total;

-- Rename and add columns in order_group_items
ALTER TABLE order_group_items RENAME COLUMN IF EXISTS total_price TO total;
ALTER TABLE order_group_items ADD COLUMN IF NOT EXISTS sub_total decimal(10,2);

-- Update sub_total for existing group items (using total as a baseline)
UPDATE order_group_items SET sub_total = total WHERE sub_total IS NULL;

