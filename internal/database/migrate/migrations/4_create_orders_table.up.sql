CREATE TABLE orders (
    order_id UUID DEFAULT gen_random_uuid(),
    user_id UUID,
    status TEXT,
    date_created TIMESTAMP DEFAULT now(),
    date_updated TIMESTAMP,

    PRIMARY KEY (order_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);