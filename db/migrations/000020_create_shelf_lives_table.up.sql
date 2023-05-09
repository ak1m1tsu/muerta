CREATE TABLE
    IF NOT EXISTS shelf_lives (
        id serial NOT NULL,
        id_product integer NOT NULL,
        id_storage integer NOT NULL,
        id_measure integer NOT NULL,
        quantity integer DEFAULT 1 NOT NULL,
        purchase_date timestamp,
        end_date timestamp NOT NULL,
        created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at timestamp,
        CONSTRAINT pk_shelf_lives PRIMARY KEY (id),
        CONSTRAINT unq_shelf_lives_id_storage UNIQUE (id_storage),
        CONSTRAINT unq_shelf_lives_id_measure UNIQUE (id_measure),
        CONSTRAINT unq_shelf_lives_id_product UNIQUE (id_product)
    );

ALTER TABLE shelf_lives
ADD
    CONSTRAINT fk_shelf_lives_products FOREIGN KEY (id_product) REFERENCES products(id);

ALTER TABLE shelf_lives
ADD
    CONSTRAINT fk_shelf_lives_storages FOREIGN KEY (id_storage) REFERENCES storages(id);

ALTER TABLE shelf_lives
ADD
    CONSTRAINT fk_shelf_lives_measures FOREIGN KEY (id_measure) REFERENCES measures(id);