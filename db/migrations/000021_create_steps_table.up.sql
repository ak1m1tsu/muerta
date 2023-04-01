CREATE TABLE
    IF NOT EXISTS steps (
        id serial NOT NULL,
        id_parent integer,
        updated_at timestamp DEFAULT CURRENT_DATE NOT NULL,
        deleted_at timestamp,
        CONSTRAINT pk_steps PRIMARY KEY (id),
        CONSTRAINT unq_steps_id_parent UNIQUE (id_parent)
    );

ALTER TABLE steps
ADD
    CONSTRAINT fk_steps_steps FOREIGN KEY (id) REFERENCES steps(id_parent);

ALTER TABLE steps
ADD
    CONSTRAINT fk_steps_recipes FOREIGN KEY (id) REFERENCES recipes(id_step);