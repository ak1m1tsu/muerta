CREATE TABLE
    IF NOT EXISTS tips (
        id serial NOT NULL,
        description text NOT NULL,
        updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at timestamp,
        CONSTRAINT pk_tips PRIMARY KEY (id)
    );