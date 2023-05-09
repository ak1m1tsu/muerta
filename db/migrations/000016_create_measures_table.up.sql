CREATE TABLE
    IF NOT EXISTS measures (
        id serial NOT NULL,
        name varchar(100) NOT NULL,
        CONSTRAINT pk_measures PRIMARY KEY (id)
    );