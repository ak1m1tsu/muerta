CREATE TABLE
    settings_categories (
        id integer DEFAULT nextval(
            'settings_categories_id_seq' :: regclass
        ) NOT NULL,
        name varchar(100) NOT NULL,
        CONSTRAINT pk_settings_categories PRIMARY KEY (id)
    );