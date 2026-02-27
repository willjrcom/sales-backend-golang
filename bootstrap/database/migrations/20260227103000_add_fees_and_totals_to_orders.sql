ALTER TABLE orders ADD COLUMN IF NOT EXISTS fees JSONB;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS sub_total decimal(10,2);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS total decimal(10,2);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS total_paid decimal(10,2);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS total_change decimal(10,2);

-- Perform the data migration
UPDATE orders SET total = total_payable WHERE total_payable IS NOT NULL AND total IS NULL;
UPDATE orders SET sub_total = total WHERE sub_total IS NULL AND total IS NOT NULL;
