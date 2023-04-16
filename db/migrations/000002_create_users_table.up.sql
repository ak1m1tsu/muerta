CREATE TABLE
    users (
        id integer DEFAULT nextval('users_id_seq' :: regclass) NOT NULL,
        name varchar(100) NOT NULL,
        salt varchar NOT NULL,
        created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at timestamp,
        CONSTRAINT pk_users PRIMARY KEY (id)
    );