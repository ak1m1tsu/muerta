CREATE TABLE
    tips (
        id integer DEFAULT nextval('tips_id_seq' :: regclass) NOT NULL,
        description text NOT NULL,
        updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at timestamp,
        CONSTRAINT pk_tips PRIMARY KEY (id)
    );