CREATE TABLE
    IF NOT EXISTS products_categories (
        id_product integer NOT NULL,
        id_category integer NOT NULL,
        CONSTRAINT pk_products_categories PRIMARY KEY (id_product, id_category),
        CONSTRAINT unq_products_categories_id_category UNIQUE (id_category),
        CONSTRAINT unq_products_categories_id_product UNIQUE (id_product)
    );

ALTER TABLE
    products_categories
ADD
    CONSTRAINT fk_products_categories_c FOREIGN KEY (id_category) REFERENCES categories(id);

ALTER TABLE
    products_categories
ADD
    CONSTRAINT fk_products_categories_p FOREIGN KEY (id_product) REFERENCES products(id);