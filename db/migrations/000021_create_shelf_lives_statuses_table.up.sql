CREATE TABLE
    IF NOT EXISTS shelf_lives_statuses (
        id_shelf_life integer NOT NULL,
        id_status integer NOT NULL,
        CONSTRAINT pk_shelf_lives_statuses PRIMARY KEY (id_shelf_life, id_status),
        CONSTRAINT unq_shelf_lives_statuses_id_shelf_life UNIQUE (id_shelf_life),
        CONSTRAINT unq_shelf_lives_statuses_id_status UNIQUE (id_status)
    );

ALTER TABLE
    shelf_lives_statuses
ADD
    CONSTRAINT fk_shelf_lives_statuses_s FOREIGN KEY (id_status) REFERENCES statuses(id);

ALTER TABLE
    shelf_lives_statuses
ADD
    CONSTRAINT fk_shelf_lives_statuses_sl FOREIGN KEY (id_shelf_life) REFERENCES shelf_lives(id);