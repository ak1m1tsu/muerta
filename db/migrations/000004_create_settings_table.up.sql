CREATE TABLE
    settings (
        id integer DEFAULT nextval('settings_id_seq' :: regclass) NOT NULL,
        id_category integer NOT NULL,
        updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at timestamp,
        CONSTRAINT pk_settings PRIMARY KEY (id),
        CONSTRAINT unq_settings_id_category UNIQUE (id_category)
    );

ALTER TABLE settings
ADD
    CONSTRAINT fk_settings FOREIGN KEY (id_category) REFERENCES settings_categories(id);