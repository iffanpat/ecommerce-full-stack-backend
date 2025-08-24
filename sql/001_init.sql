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

-- Seed Products Data
INSERT INTO products (sku, name, description, price_cents, stock)
VALUES
 ('SKU-001','Basic Tee','Cotton tee', 19900, 50),
 ('SKU-002','Hoodie','Cozy hoodie', 59900, 20),
 ('SKU-003','Cap','Classic cap', 14900, 35),
 ('SKU-004','Jeans','Denim jeans', 79900, 15),
 ('SKU-005','Sneakers','Casual sneakers', 129900, 25)
ON CONFLICT (sku) DO NOTHING;

-- Seed Product Images
INSERT INTO product_images (product_id, url, is_primary) VALUES
 (1, 'https://picsum.photos/seed/tee/600', TRUE),
 (2, 'https://picsum.photos/seed/hoodie/600', TRUE),
 (3, 'https://picsum.photos/seed/cap/600', TRUE),
 (4, 'https://picsum.photos/seed/jeans/600', TRUE),
 (5, 'https://picsum.photos/seed/sneakers/600', TRUE)
ON CONFLICT DO NOTHING;
      