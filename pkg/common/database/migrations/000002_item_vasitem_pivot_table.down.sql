DROP TABLE IF EXISTS item_vas_items;

ALTER TABLE IF EXISTS vas_items
    ADD COLUMN IF NOT EXISTS item_id int;