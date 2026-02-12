ALTER TABLE products RENAME COLUMN code TO sku;
ALTER TABLE stock_alerts RENAME COLUMN product_code TO product_sku;
