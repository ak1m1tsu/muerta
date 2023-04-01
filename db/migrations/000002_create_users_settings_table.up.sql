CREATE TABLE
    IF NOT EXISTS users_settings (
        id_user integer NOT NULL,
        id_setting integer NOT NULL,
        "value" varchar NOT NULL,
        CONSTRAINT pk_users_settings PRIMARY KEY (id_user, id_setting),
        CONSTRAINT unq_users_settings_id_user UNIQUE (id_user),
        CONSTRAINT unq_users_settings_id_setting UNIQUE (id_setting)
    );