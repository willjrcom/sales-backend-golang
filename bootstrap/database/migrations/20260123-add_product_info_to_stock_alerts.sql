-- Add product_name and product_sku to stock_alerts table

-- Up
ALTER TABLE stock_alerts ADD COLUMN IF NOT EXISTS product_name TEXT;
ALTER TABLE stock_alerts ADD COLUMN IF NOT EXISTS product_code TEXT;
