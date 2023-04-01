CREATE TABLE
    IF NOT EXISTS storages (
        id serial NOT NULL,
        id_type integer NOT NULL,
        name varchar(100) NOT NULL,
        temperature real,
        humidity real,
        created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at timestamp,
        CONSTRAINT pk_storages PRIMARY KEY (id),
        CONSTRAINT unq_storages_id_type UNIQUE (id_type)
    );

ALTER TABLE storages
ADD
    CONSTRAINT fk_storages_users_storages FOREIGN KEY (id) REFERENCES users_storages(id_storage);

ALTER TABLE storages
ADD
    CONSTRAINT fk_storages_storages_tips FOREIGN KEY (id) REFERENCES storages_tips(id_storage);

ALTER TABLE storages
ADD
    CONSTRAINT fk_storages_shelf_lives FOREIGN KEY (id) REFERENCES shelf_lives(id_storage);