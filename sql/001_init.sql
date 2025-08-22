-- USERS
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE,
  password_hash TEXT,
  created_at TIMESTAMP DEFAULT now()
);

-- PRODUCTS
CREATE TABLE IF NOT EXISTS products (
  id SERIAL PRIMARY KEY,
  sku TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  price_cents INT NOT NULL,
  currency TEXT NOT NULL DEFAULT 'THB',
  stock INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS product_images (
  id SERIAL PRIMARY KEY,
  product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  url TEXT NOT NULL,
  is_primary BOOLEAN DEFAULT FALSE
);

-- CARTS
CREATE TABLE IF NOT EXISTS carts (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id),
  guest_token TEXT,
  created_at TIMESTAMP DEFAULT now()
);
-- Partial unique index for guest_token (only when not null)
CREATE UNIQUE INDEX IF NOT EXISTS ux_carts_guest ON carts(guest_token) WHERE guest_token IS NOT NULL;

CREATE TABLE IF NOT EXISTS cart_items (
  id SERIAL PRIMARY KEY,
  cart_id INT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
  product_id INT NOT NULL REFERENCES products(id),
  qty INT NOT NULL CHECK (qty > 0),
  UNIQUE(cart_id, product_id)
);

-- ORDERS
CREATE TABLE IF NOT EXISTS orders (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id),
  total_cents INT NOT NULL,
  currency TEXT NOT NULL DEFAULT 'THB',
  status TEXT NOT NULL DEFAULT 'PENDING',
  created_at TIMESTAMP DEFAULT now(),
  cart_id INT
);

CREATE TABLE IF NOT EXISTS order_items (
  id SERIAL PRIMARY KEY,
  order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id INT NOT NULL REFERENCES products(id),
  price_cents INT NOT NULL,
  qty INT NOT NULL
);

-- Idempotency
CREATE TABLE IF NOT EXISTS idempotency_keys (
  id SERIAL PRIMARY KEY,
  key TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT now()
);

-- Seed
INSERT INTO products (sku, name, description, price_cents, stock)
VALUES
 ('SKU-001','Basic Tee','Cotton tee', 19900, 50),
 ('SKU-002','Hoodie','Cozy hoodie', 59900, 20),
 ('SKU-003','Cap','Classic cap', 14900, 35)
ON CONFLICT DO NOTHING;

INSERT INTO product_images (product_id, url, is_primary) VALUES
 (1, 'https://picsum.photos/seed/tee/600', TRUE),
 (2, 'https://picsum.photos/seed/hoodie/600', TRUE),
 (3, 'https://picsum.photos/seed/cap/600', TRUE)
ON CONFLICT DO NOTHING;
      