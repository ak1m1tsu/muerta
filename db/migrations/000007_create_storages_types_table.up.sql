CREATE TABLE
    storages_types (
        id integer DEFAULT nextval(
            'storages_types_id_seq' :: regclass
        ) NOT NULL,
        name varchar(100) NOT NULL,
        CONSTRAINT pk_storages_types PRIMARY KEY (id)
    );