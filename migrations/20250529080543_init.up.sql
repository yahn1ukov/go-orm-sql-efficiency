CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,
    stock INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    total NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_products (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);

INSERT INTO customers (name, email) VALUES
('John Doe', 'john.doe@example.com'),
('Jane Smith', 'jane.smith@example.com'),
('Bob Johnson', 'bob.johnson@example.com'),
('Alice Williams', 'alice.williams@example.com');

INSERT INTO products (name, description, price, stock) VALUES
('Laptop', 'High-performance laptop', 1299.99, 50),
('Smartphone', 'Latest smartphone model', 799.99, 100),
('Tablet', '10-inch display tablet', 499.99, 75),
('Headphones', 'Wireless noise-cancelling', 199.99, 200);

INSERT INTO orders (customer_id, total) VALUES
(1, 1999.98),
(2, 1199.97),
(3, 899.98),
(4, 2999.95);

INSERT INTO order_products (order_id, product_id, quantity, price) VALUES
(1, 1, 1, 1299.99),
(1, 4, 2, 199.99),
(2, 2, 1, 799.99),
(2, 4, 1, 199.99),
(3, 3, 1, 499.99),
(3, 4, 2, 199.99),
(4, 1, 2, 1299.99),
(4, 2, 1, 799.99),
(4, 3, 1, 499.99);
