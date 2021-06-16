-- Many to many relationship table
-- Many orders can be related to one or many products
-- Many products can be related to one or many orders
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