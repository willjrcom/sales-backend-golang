-- =============================================================================
-- Migration consolidada: estoque com lotes, reservas e alertas aprimorados
-- Data: 2026-03-02
-- =============================================================================

-- 1. Renomear colunas em order_items (sub_total e total)
ALTER TABLE IF EXISTS order_items RENAME COLUMN price TO sub_total;
ALTER TABLE IF EXISTS order_items RENAME COLUMN total_price TO total;

-- 2. Renomear coluna em order_group_items + adicionar sub_total
ALTER TABLE IF EXISTS order_group_items RENAME COLUMN total_price TO total;
ALTER TABLE IF EXISTS order_group_items ADD COLUMN IF NOT EXISTS sub_total DECIMAL(10,2);
UPDATE order_group_items SET sub_total = total WHERE sub_total IS NULL;

-- 3. Remover coluna cost de product_variations
ALTER TABLE IF EXISTS product_variations DROP COLUMN IF EXISTS cost;

-- 4. Adicionar reserved_stock em stocks
ALTER TABLE IF EXISTS stocks
    ADD COLUMN IF NOT EXISTS reserved_stock DECIMAL(10,3) NOT NULL DEFAULT 0;

-- 5. Permitir product_variation_id nulo em stock_alerts e stock_batches
--    (produtos sem variação não têm variation_id válido)
ALTER TABLE IF EXISTS stock_alerts  ALTER COLUMN product_variation_id DROP NOT NULL;
ALTER TABLE IF EXISTS stock_batches ALTER COLUMN product_variation_id DROP NOT NULL;
