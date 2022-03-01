DROP TABLE IF EXISTS urls CASCADE;

CREATE TABLE urls
(
    id         SERIAL,
    org_url    varchar(255) NOT NULL UNIQUE,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);