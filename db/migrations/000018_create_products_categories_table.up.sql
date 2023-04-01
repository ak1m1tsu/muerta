CREATE TABLE
    IF NOT EXISTS products_categories (
        id_product integer NOT NULL,
        id_category integer NOT NULL,
        CONSTRAINT pk_products_categories PRIMARY KEY (id_product, id_category),
        CONSTRAINT unq_products_categories_id_category UNIQUE (id_category),
        CONSTRAINT unq_products_categories_id_product UNIQUE (id_product)
    );