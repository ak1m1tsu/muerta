CREATE TABLE
    IF NOT EXISTS measures (
        id serial NOT NULL,
        name varchar(100) NOT NULL,
        CONSTRAINT pk_measures PRIMARY KEY (id)
    );

ALTER TABLE measures
ADD
    CONSTRAINT fk_measures_shelf_lives FOREIGN KEY (id) REFERENCES shelf_lives(id_measure);

ALTER TABLE measures
ADD
    CONSTRAINT fk_measures FOREIGN KEY (id) REFERENCES products_recipes_measures(id_measure);