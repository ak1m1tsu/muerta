CREATE TABLE
    steps (
        id integer DEFAULT nextval('steps_id_seq' :: regclass) NOT NULL,
        id_parent integer,
        updated_at timestamp DEFAULT CURRENT_DATE NOT NULL,
        deleted_at timestamp,
        CONSTRAINT pk_steps PRIMARY KEY (id),
        CONSTRAINT unq_steps_id_parent UNIQUE (id_parent)
    );

ALTER TABLE steps
ADD
    CONSTRAINT fk_steps_steps FOREIGN KEY (id_parent) REFERENCES steps(id);