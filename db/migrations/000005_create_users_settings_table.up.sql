CREATE TABLE
    IF NOT EXISTS users_settings (
        id_user integer NOT NULL,
        id_setting integer NOT NULL,
        "value" varchar NOT NULL,
        CONSTRAINT unq_users_settings_id_user UNIQUE (id_user),
        CONSTRAINT unq_users_settings_id_setting UNIQUE (id_setting),
        CONSTRAINT pk_users_settings PRIMARY KEY (id_user, id_setting)
    );

ALTER TABLE users_settings
ADD
    CONSTRAINT fk_users_settings_users FOREIGN KEY (id_user) REFERENCES users(id);

ALTER TABLE users_settings
ADD
    CONSTRAINT fk_users_settings_settings FOREIGN KEY (id_setting) REFERENCES settings(id);