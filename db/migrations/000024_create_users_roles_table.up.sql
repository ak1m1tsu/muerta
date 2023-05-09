CREATE TABLE
    IF NOT EXISTS users_roles (
        id_user integer NOT NULL,
        id_role integer NOT NULL,
        CONSTRAINT pk_users_roles PRIMARY KEY (id_user, id_role)
    );

ALTER TABLE users_roles
ADD
    CONSTRAINT fk_users_roles_users FOREIGN KEY (id_user) REFERENCES users(id);

ALTER TABLE users_roles
ADD
    CONSTRAINT fk_users_roles_roles FOREIGN KEY (id_role) REFERENCES roles(id);