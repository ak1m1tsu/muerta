CREATE TABLE
    IF NOT EXISTS steps (
        id serial NOT NULL,
        updated_at timestamp DEFAULT CURRENT_DATE NOT NULL,
        deleted_at timestamp,
        name varchar(100) NOT NULL,
        CONSTRAINT pk_steps PRIMARY KEY (id)
    );