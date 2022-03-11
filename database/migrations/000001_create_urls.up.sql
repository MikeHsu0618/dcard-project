DROP TABLE IF EXISTS urls CASCADE;

CREATE TABLE urls
(
    id         BIGSERIAL,
    org_url    varchar(255) NOT NULL UNIQUE,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

COMMENT ON COLUMN urls.org_url IS '原網址'