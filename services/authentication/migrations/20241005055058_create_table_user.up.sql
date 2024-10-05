CREATE TABLE "user"
(
    id          UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name  VARCHAR(256)    NOT NULL,
    last_name   VARCHAR(256),
    username    VARCHAR(256)    NOT NULL UNIQUE,
    email       VARCHAR(256)    NOT NULL UNIQUE,
    salt        VARCHAR(256)    NOT NULL,
    password    VARCHAR(256)    NOT NULL,
    created_by  VARCHAR(256),
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(256),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_deleted   BOOLEAN DEFAULT FALSE
);