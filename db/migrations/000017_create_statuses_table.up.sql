CREATE TABLE
    IF NOT EXISTS statuses (
        id serial NOT NULL,
        name varchar(100) NOT NULL,
        CONSTRAINT pk_statuses PRIMARY KEY (id)
    );