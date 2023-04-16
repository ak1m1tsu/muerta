CREATE TABLE
    statuses (
        id integer DEFAULT nextval('statuses_id_seq' :: regclass) NOT NULL,
        name varchar(100) NOT NULL,
        CONSTRAINT pk_statuses PRIMARY KEY (id)
    );