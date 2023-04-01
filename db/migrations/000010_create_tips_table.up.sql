CREATE TABLE
    IF NOT EXISTS tips (
        id serial NOT NULL,
        description text NOT NULL,
        updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at timestamp,
        CONSTRAINT pk_tips PRIMARY KEY (id)
    );

ALTER TABLE tips
ADD
    CONSTRAINT fk_tips_types_storages_tips FOREIGN KEY (id) REFERENCES storages_types_tips(id_tip);

ALTER TABLE tips
ADD
    CONSTRAINT fk_tips_storages_tips FOREIGN KEY (id) REFERENCES storages_tips(id_tip);

ALTER TABLE tips
ADD
    CONSTRAINT fk_tips_products_tips FOREIGN KEY (id) REFERENCES products_tips(id_tip);