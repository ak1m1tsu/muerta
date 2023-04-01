CREATE TABLE
    IF NOT EXISTS settings_categories (
        id serial NOT NULL,
        name varchar(100) NOT NULL,
        CONSTRAINT pk_settings_categories PRIMARY KEY (id)
    );

ALTER TABLE
    settings_categories
ADD
    CONSTRAINT fk_settings_categories FOREIGN KEY (id) REFERENCES settings(id_category);