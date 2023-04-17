CREATE TABLE
    recipes (
        id integer DEFAULT nextval('recipes_id_seq' :: regclass) NOT NULL,
        id_step integer NOT NULL,
        name varchar(100) NOT NULL,
        description text,
        updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at timestamp,
        CONSTRAINT pk_recipes PRIMARY KEY (id),
        CONSTRAINT unq_recipes_id_step UNIQUE (id_step)
    );

ALTER TABLE recipes
ADD
    CONSTRAINT fk_recipes_steps FOREIGN KEY (id_step) REFERENCES steps(id);