CREATE TABLE if not exists "product_variations" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "product_id" uuid NOT NULL,
  "size_id" uuid NOT NULL,
  "price" numeric(10,2) NOT NULL,
  "cost" numeric(10,2),
  "is_available" boolean DEFAULT true,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now(),
  "deleted_at" timestamptz,
  PRIMARY KEY ("id"),
);

-- Drop columns from products
ALTER TABLE "products" DROP COLUMN IF EXISTS "price";
ALTER TABLE "products" DROP COLUMN IF EXISTS "cost";
ALTER TABLE "products" DROP COLUMN IF EXISTS "is_available";
ALTER TABLE "products" DROP COLUMN IF EXISTS "size_id";
