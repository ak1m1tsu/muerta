CREATE TABLE
    IF NOT EXISTS settings_categories (
        id serial NOT NULL,
        name varchar(100) NOT NULL,
        CONSTRAINT pk_settings_categories PRIMARY KEY (id)
    );