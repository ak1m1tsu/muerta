CREATE TABLE
    measures (
        id integer DEFAULT nextval('measures_id_seq' :: regclass) NOT NULL,
        name varchar(100) NOT NULL,
        CONSTRAINT pk_measures PRIMARY KEY (id)
    );