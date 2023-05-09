CREATE TABLE
    IF NOT EXISTS users_storages (
        id_user integer NOT NULL,
        id_storage integer NOT NULL,
        CONSTRAINT unq_users_storages_id_storage UNIQUE (id_storage),
        CONSTRAINT unq_users_storages_id_user UNIQUE (id_user),
        CONSTRAINT pk_users_storages PRIMARY KEY (id_user, id_storage)
    );

ALTER TABLE users_storages
ADD
    CONSTRAINT fk_users_storages_users FOREIGN KEY (id_user) REFERENCES users(id);

ALTER TABLE users_storages
ADD
    CONSTRAINT fk_users_storages_storages FOREIGN KEY (id_storage) REFERENCES storages(id);