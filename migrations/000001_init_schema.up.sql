
BEGIN;

CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash TEXT NOT NULL,
                       display_name VARCHAR(120) NOT NULL DEFAULT '',
                       role VARCHAR(20) NOT NULL CHECK (role IN ('customer', 'admin')),
                       is_active BOOLEAN NOT NULL DEFAULT TRUE,
                       created_at TIMESTAMP NOT NULL DEFAULT now(),
                       updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE categories (
                            id BIGSERIAL PRIMARY KEY,
                            name VARCHAR(120) NOT NULL,
                            slug VARCHAR(140) NOT NULL UNIQUE,
                            parent_id BIGINT REFERENCES categories (id) ON DELETE SET NULL,
                            sort_order INT NOT NULL DEFAULT 0 CHECK (sort_order >= 0),
                            is_active BOOLEAN NOT NULL DEFAULT TRUE,
                            created_at TIMESTAMP NOT NULL DEFAULT now(),
                            updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE products (
                          id BIGSERIAL PRIMARY KEY,
                          category_id BIGINT NOT NULL REFERENCES categories (id) ON DELETE RESTRICT,
                          sku VARCHAR(64) NOT NULL UNIQUE,
                          name VARCHAR(255) NOT NULL,
                          description TEXT,
                          price NUMERIC(12, 2) NOT NULL CHECK (price > 0),
                          stock INT NOT NULL CHECK (stock >= 0),
                          is_active BOOLEAN NOT NULL DEFAULT TRUE,
                          created_at TIMESTAMP NOT NULL DEFAULT now(),
                          updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE product_images (
                                id BIGSERIAL PRIMARY KEY,
                                product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
                                url TEXT NOT NULL,
                                sort_order INT NOT NULL DEFAULT 0 CHECK (sort_order >= 0),
                                created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE addresses (
                           id BIGSERIAL PRIMARY KEY,
                           user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
                           city VARCHAR(120) NOT NULL,
                           street VARCHAR(160) NOT NULL,
                           house VARCHAR(32) NOT NULL,
                           apartment VARCHAR(32),
                           postal_code VARCHAR(20),
                           comment TEXT,
                           created_at TIMESTAMP NOT NULL DEFAULT now(),
                           updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE carts (
                       id BIGSERIAL PRIMARY KEY,
                       user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
                       status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'ordered', 'abandoned')),
                       created_at TIMESTAMP NOT NULL DEFAULT now(),
                       updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ux_carts_one_active_per_user
    ON carts (user_id)
    WHERE status = 'active';

CREATE TABLE cart_items (
                            id BIGSERIAL PRIMARY KEY,
                            cart_id BIGINT NOT NULL REFERENCES carts (id) ON DELETE CASCADE,
                            product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE RESTRICT,
                            qty INT NOT NULL CHECK (qty > 0),
                            price_snapshot NUMERIC(12, 2) NOT NULL CHECK (price_snapshot > 0),
                            created_at TIMESTAMP NOT NULL DEFAULT now(),
                            updated_at TIMESTAMP NOT NULL DEFAULT now(),
                            UNIQUE (cart_id, product_id)
);

CREATE TABLE orders (
                        id BIGSERIAL PRIMARY KEY,
                        user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
                        status VARCHAR(20) NOT NULL CHECK (status IN ('new', 'paid', 'shipped', 'done', 'cancelled')),
                        total_price NUMERIC(12, 2) NOT NULL CHECK (total_price >= 0),
                        address_id BIGINT REFERENCES addresses (id) ON DELETE SET NULL,
                        comment TEXT,
                        created_at TIMESTAMP NOT NULL DEFAULT now(),
                        updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE order_items (
                             id BIGSERIAL PRIMARY KEY,
                             order_id BIGINT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
                             product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE RESTRICT,
                             qty INT NOT NULL CHECK (qty > 0),
                             unit_price NUMERIC(12, 2) NOT NULL CHECK (unit_price > 0),
                             line_total NUMERIC(12, 2) NOT NULL CHECK (line_total >= 0),
                             CHECK (line_total = (qty::NUMERIC * unit_price))
);

CREATE TABLE audit_logs (
                            id BIGSERIAL PRIMARY KEY,
                            actor_id BIGINT REFERENCES users (id) ON DELETE SET NULL,
                            entity VARCHAR(80) NOT NULL,
                            entity_id BIGINT,
                            action VARCHAR(80) NOT NULL CHECK (action IN ('create', 'update', 'delete', 'status_change')),
    payload JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_products_category_active ON products (category_id, is_active);
CREATE INDEX idx_products_price ON products (price);
CREATE INDEX idx_categories_parent ON categories (parent_id);
CREATE INDEX idx_orders_user_status_created ON orders (user_id, status, created_at);
CREATE INDEX idx_order_items_order ON order_items (order_id);
CREATE INDEX idx_cart_items_cart ON cart_items (cart_id);
CREATE INDEX idx_audit_logs_entity ON audit_logs (entity, entity_id);

COMMIT;