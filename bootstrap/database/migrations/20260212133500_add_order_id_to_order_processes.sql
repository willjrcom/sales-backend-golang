ALTER TABLE order_processes ADD COLUMN order_id UUID;
UPDATE order_processes op SET order_id = gi.order_id FROM order_group_items gi WHERE op.group_item_id = gi.id;
ALTER TABLE order_processes ALTER COLUMN order_id SET NOT NULL;
