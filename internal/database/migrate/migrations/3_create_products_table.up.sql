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