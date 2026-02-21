-- Remove duplicate active processes that may exist before applying the constraint
-- Keep only the oldest one per (group_item_id, process_rule_id) combination
DELETE FROM order_processes
WHERE id NOT IN (
    SELECT DISTINCT ON (group_item_id, process_rule_id) id
    FROM order_processes
    WHERE finished_at IS NULL
    ORDER BY group_item_id, process_rule_id, created_at ASC
)
AND finished_at IS NULL;

-- Unique index: only one active (not finished) process per group_item + process_rule
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_active_process
    ON order_processes (group_item_id, process_rule_id)
    WHERE finished_at IS NULL;
