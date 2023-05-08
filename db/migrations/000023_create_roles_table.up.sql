CREATE TABLE
    IF NOT EXISTS roles (
        id serial NOT NULL,
        name varchar(100) NOT NULL,
        created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at timestamptz,
        CONSTRAINT pk_roles PRIMARY KEY (id)
    );