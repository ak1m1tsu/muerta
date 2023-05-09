CREATE TABLE
    IF NOT EXISTS recipes (
        id serial NOT NULL,
        id_user integer NOT NULL,
        name varchar(100) NOT NULL,
        description text,
        updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at timestamp,
        CONSTRAINT pk_recipes PRIMARY KEY (id),
        CONSTRAINT unq_recipes_name UNIQUE (name)
    );

ALTER TABLE recipes
ADD
    CONSTRAINT fk_recipes_users FOREIGN KEY (id_user) REFERENCES users(id);