CREATE TABLE
    IF NOT EXISTS products_tips (
        id_product integer NOT NULL,
        id_tip integer NOT NULL,
        CONSTRAINT pk_products_tips PRIMARY KEY (id_product, id_tip),
        CONSTRAINT unq_products_tips_id_tip UNIQUE (id_tip),
        CONSTRAINT unq_products_tips_id_product UNIQUE (id_product)
    );