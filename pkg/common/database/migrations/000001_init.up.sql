CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    item_id INT,
    category_id INT,
    seller_id INT,
    price DECIMAL(10, 2),
    quantity INT
    );

CREATE TABLE IF NOT EXISTS vas_items (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    item_id INT,
    vas_item_id INT,
    category_id INT,
    seller_id INT,
    price DECIMAL(10, 2),
    quantity INT
    );
