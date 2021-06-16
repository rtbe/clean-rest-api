
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