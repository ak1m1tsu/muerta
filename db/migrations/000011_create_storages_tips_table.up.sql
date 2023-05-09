CREATE TABLE
    IF NOT EXISTS storages_tips (
        id_storage integer NOT NULL,
        id_tip integer NOT NULL,
        CONSTRAINT pk_storages_tips PRIMARY KEY (id_storage, id_tip),
        CONSTRAINT unq_storages_tips_id_storage UNIQUE (id_storage),
        CONSTRAINT unq_storages_tips_id_tip UNIQUE (id_tip)
    );

ALTER TABLE storages_tips
ADD
    CONSTRAINT fk_storages_tips_storages FOREIGN KEY (id_storage) REFERENCES storages(id);

ALTER TABLE storages_tips
ADD
    CONSTRAINT fk_storages_tips_tips FOREIGN KEY (id_tip) REFERENCES tips(id);