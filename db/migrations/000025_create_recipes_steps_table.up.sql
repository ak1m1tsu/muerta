CREATE TABLE
    IF NOT EXISTS recipes_steps (
        id_recipe integer NOT NULL,
        id_step integer NOT NULL,
        place integer NOT NULL,
        CONSTRAINT pk_recipes_steps PRIMARY KEY (id_recipe, id_step, place)
    );

ALTER TABLE recipes_steps
ADD
    CONSTRAINT fk_recipes_steps_recipes FOREIGN KEY (id_recipe) REFERENCES recipes(id);

ALTER TABLE recipes_steps
ADD
    CONSTRAINT fk_recipes_steps_steps FOREIGN KEY (id_step) REFERENCES steps(id);