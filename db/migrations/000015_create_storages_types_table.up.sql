CREATE TABLE
    IF NOT EXISTS storages_types (
        id serial NOT NULL,
        name varchar(100) NOT NULL,
        CONSTRAINT pk_storages_types PRIMARY KEY (id)
    );

ALTER TABLE storages_types
ADD
    CONSTRAINT fk_storages_types_storages FOREIGN KEY (id) REFERENCES storages(id_type);

ALTER TABLE storages_types
ADD
    CONSTRAINT fk_storages_types FOREIGN KEY (id) REFERENCES storages_types_tips(id_storage_type);