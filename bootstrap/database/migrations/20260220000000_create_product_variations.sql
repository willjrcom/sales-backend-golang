CREATE TABLE "product_variations" (
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
  FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  FOREIGN KEY ("size_id") REFERENCES "sizes" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Migrate existing data
INSERT INTO "product_variations" ("product_id", "size_id", "price", "cost", "is_available", "created_at", "updated_at", "deleted_at")
SELECT "id", "size_id", "price", "cost", "is_available", "created_at", "updated_at", "deleted_at"
FROM "products";

-- Drop columns from products
ALTER TABLE "products" DROP COLUMN "price";
ALTER TABLE "products" DROP COLUMN "cost";
ALTER TABLE "products" DROP COLUMN "is_available";
ALTER TABLE "products" DROP COLUMN "size_id";
