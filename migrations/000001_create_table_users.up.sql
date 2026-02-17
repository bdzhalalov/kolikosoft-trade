CREATE TABLE users (
       id           BIGSERIAL PRIMARY KEY,
       balance      BIGINT NOT NULL DEFAULT 0,
       created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
       updated_at   TIMESTAMP NULL
);

CREATE OR REPLACE FUNCTION trigger_set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();
