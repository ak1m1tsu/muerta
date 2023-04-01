CREATE TABLE
    IF NOT EXISTS users (
        id serial NOT NULL,
        name varchar(100) NOT NULL,
        salt varchar NOT NULL,
        created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at timestamp,
        CONSTRAINT pk_users PRIMARY KEY (id)
    );

ALTER TABLE users
ADD
    CONSTRAINT fk_users_users_storages FOREIGN KEY (id) REFERENCES users_storages(id_user);

ALTER TABLE users
ADD
    CONSTRAINT fk_users_users_settings FOREIGN KEY (id) REFERENCES users_settings(id_user);