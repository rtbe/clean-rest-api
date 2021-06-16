INSERT INTO users  
    (user_name, first_name, last_name, password, email, roles)
VALUES 
    ('admin', 'admin', 'admin', '$2y$12$IdauwDQYA4dvOc8tvfXu5eCcw8tAkNPZ8w2FTHxf2jcPAOYCzfozi', 'admin@example.com', ARRAY ['ADMIN']);