CREATE TABLE
    IF NOT EXISTS users_storages (
        id_user integer NOT NULL,
        id_storage integer NOT NULL,
        CONSTRAINT pk_users_storages PRIMARY KEY (id_user, id_storage),
        CONSTRAINT unq_users_storages_id_storage UNIQUE (id_storage),
        CONSTRAINT unq_users_storages_id_user UNIQUE (id_user)
    );