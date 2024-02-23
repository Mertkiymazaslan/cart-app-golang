CREATE TABLE IF NOT EXISTS item_vas_items (
      id SERIAL,
      created_at TIMESTAMPTZ,
      updated_at TIMESTAMPTZ,
      deleted_at TIMESTAMPTZ,
      item_id INT,
      vas_item_id INT
);

ALTER TABLE IF EXISTS vas_items
    DROP COLUMN IF EXISTS item_id;