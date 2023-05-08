CREATE TABLE
    IF NOT EXISTS steps (
        id integer DEFAULT nextval('steps_id_seq' :: regclass) NOT NULL,
        updated_at timestamp DEFAULT CURRENT_DATE NOT NULL,
        deleted_at timestamp,
        name varchar(100) NOT NULL,
        CONSTRAINT pk_steps PRIMARY KEY (id)
    );