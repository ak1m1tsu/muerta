CREATE TABLE
    IF NOT EXISTS storages_types (
        id serial NOT NULL,
        name varchar(100) NOT NULL,
        CONSTRAINT pk_storages_types PRIMARY KEY (id)
    );