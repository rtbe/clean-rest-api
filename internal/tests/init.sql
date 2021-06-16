-- Set of initialisation commands to prepare a database for running tests.
CREATE EXTENSION pgcrypto;

CREATE TABLE users (
    user_id UUID DEFAULT gen_random_uuid(),
    user_name TEXT UNIQUE NOT NULL,
    first_name TEXT, 
    last_name TEXT, 
    password TEXT, 
    email TEXT UNIQUE NOT NULL, 
    roles TEXT[],
    date_created TIMESTAMP DEFAULT now(), 
    date_updated TIMESTAMP, 

    PRIMARY KEY (user_id)
);

CREATE TABLE products (
    product_id UUID DEFAULT gen_random_uuid(), 
    title TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10,2),
    stock INT,
    date_created TIMESTAMP DEFAULT now(),
    date_updated TIMESTAMP,

    PRIMARY KEY (product_id)
);

CREATE TABLE orders (
    order_id UUID DEFAULT gen_random_uuid(),
    user_id UUID,
    status TEXT,
    date_created TIMESTAMP DEFAULT now(),
    date_updated TIMESTAMP,

    PRIMARY KEY (order_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

CREATE TABLE order_items (
    order_item_id UUID DEFAULT gen_random_uuid(),
    order_id UUID,
    product_id UUID,
    quantity INT,
    date_created TIMESTAMP DEFAULT now(),
    date_updated TIMESTAMP,
    
    PRIMARY KEY (order_item_id),
    FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products (product_id) ON DELETE CASCADE
);

CREATE INDEX idx_order_to_product ON order_items (order_id, product_id);