CREATE TABLE
    IF NOT EXISTS storages_types_tips (
        id_tip integer NOT NULL,
        id_storage_type integer NOT NULL,
        CONSTRAINT pk_types_storages_tips PRIMARY KEY (id_tip, id_storage_type),
        CONSTRAINT unq_types_storages_tips_id_storage_type UNIQUE (id_storage_type),
        CONSTRAINT unq_types_storages_tips_id_tip UNIQUE (id_tip)
    );

ALTER TABLE
    storages_types_tips
ADD
    CONSTRAINT fk_storages_types_tips FOREIGN KEY (id_storage_type) REFERENCES storages_types(id);

ALTER TABLE
    storages_types_tips
ADD
    CONSTRAINT fk_storages_types_tips_tips FOREIGN KEY (id_tip) REFERENCES tips(id);