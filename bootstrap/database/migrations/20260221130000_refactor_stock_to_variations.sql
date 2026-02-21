-- Add product_variation_id to order_items
ALTER TABLE "order_items" ADD COLUMN IF NOT EXISTS "product_variation_id" uuid;
ALTER TABLE "order_items" ADD CONSTRAINT fk_order_items_product_variation FOREIGN KEY ("product_variation_id") REFERENCES "product_variations" ("id") ON DELETE SET NULL;

-- Add product_variation_id to stocks
ALTER TABLE "stocks" ADD COLUMN IF NOT EXISTS "product_variation_id" uuid;
ALTER TABLE "stocks" ADD CONSTRAINT fk_stocks_product_variation FOREIGN KEY ("product_variation_id") REFERENCES "product_variations" ("id") ON DELETE SET NULL;

-- Add product_id and product_variation_id to stock_alerts
ALTER TABLE "stock_alerts" ADD COLUMN IF NOT EXISTS "product_id" uuid;
ALTER TABLE "stock_alerts" ADD COLUMN IF NOT EXISTS "product_variation_id" uuid;
ALTER TABLE "stock_alerts" ADD CONSTRAINT fk_stock_alerts_product FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE;
ALTER TABLE "stock_alerts" ADD CONSTRAINT fk_stock_alerts_product_variation FOREIGN KEY ("product_variation_id") REFERENCES "product_variations" ("id") ON DELETE CASCADE;

-- Data migration: Populate variations for existing items and stocks
UPDATE "order_items" oi
SET "product_variation_id" = (SELECT id FROM "product_variations" pv WHERE pv.product_id = oi.product_id LIMIT 1)
WHERE "product_variation_id" IS NULL AND "product_id" IS NOT NULL;

UPDATE "stocks" s
SET "product_variation_id" = (SELECT id FROM "product_variations" pv WHERE pv.product_id = s.product_id LIMIT 1)
WHERE "product_variation_id" IS NULL AND "product_id" IS NOT NULL;

-- Data migration: Link alerts to product and variation
UPDATE "stock_alerts" sa
SET "product_id" = s.product_id,
    "product_variation_id" = s.product_variation_id
FROM "stocks" s
WHERE sa.stock_id = s.id AND sa.product_id IS NULL;

-- Drop obsolete text columns from stock_alerts if they exist
ALTER TABLE "stock_alerts" DROP COLUMN IF EXISTS "product_name";
ALTER TABLE "stock_alerts" DROP COLUMN IF EXISTS "product_sku";
