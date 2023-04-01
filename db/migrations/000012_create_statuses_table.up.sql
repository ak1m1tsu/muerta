CREATE TABLE
    IF NOT EXISTS statuses (
        id serial NOT NULL,
        name varchar(100) NOT NULL,
        CONSTRAINT pk_statuses PRIMARY KEY (id)
    );

ALTER TABLE statuses
ADD
    CONSTRAINT fk_statuses FOREIGN KEY (id) REFERENCES shelf_lives_statuses(id_status);