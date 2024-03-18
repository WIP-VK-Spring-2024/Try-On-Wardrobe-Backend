
-- +migrate Up
CREATE TYPE gender AS ENUM ('male', 'female', 'unisex', 'unknown');
CREATE TYPE season AS ENUM ('winter', 'spring', 'summer', 'autumn');

-- +migrate Down
DROP TYPE gender;
DROP TYPE season;
