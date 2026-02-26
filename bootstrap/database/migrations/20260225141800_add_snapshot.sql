-- Migration: Create order_group_item_snapshots table
CREATE TABLE IF NOT EXISTS order_group_item_snapshots (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    group_item_id UUID NOT NULL UNIQUE REFERENCES order_group_items(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);
