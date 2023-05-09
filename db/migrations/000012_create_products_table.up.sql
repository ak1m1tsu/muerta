CREATE TABLE
    IF NOT EXISTS products (
        id serial NOT NULL,
        name varchar(100) NOT NULL,
        updated_at timestamp DEFAULT CURRENT_DATE NOT NULL,
        deleted_at timestamp,
        CONSTRAINT pk_products PRIMARY KEY (id)
    );