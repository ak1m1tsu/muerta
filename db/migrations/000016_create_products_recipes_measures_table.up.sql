CREATE TABLE
    IF NOT EXISTS products_recipes_measures (
        id_product integer NOT NULL,
        id_recipe integer NOT NULL,
        id_measure integer NOT NULL,
        quantity real DEFAULT 1 NOT NULL,
        CONSTRAINT pk_products_recipes_measures PRIMARY KEY (
            id_product,
            id_recipe,
            id_measure
        ),
        CONSTRAINT unq_products_recipes_measures_id_measure UNIQUE (id_measure),
        CONSTRAINT unq_products_recipes_measures_id_recipe UNIQUE (id_recipe),
        CONSTRAINT unq_products_recipes_measures_id_product UNIQUE (id_product)
    );