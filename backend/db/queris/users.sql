CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    address VARCHAR(42),
    password VARCHAR(255)
);

CREATE INDEX indx_address ON users(address);