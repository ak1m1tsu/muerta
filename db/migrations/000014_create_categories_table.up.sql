CREATE TABLE
    IF NOT EXISTS categories (
        id serial NOT NULL,
        name varchar(100) NOT NULL,
        created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at timestamp,
        CONSTRAINT pk_categories PRIMARY KEY (id)
    );