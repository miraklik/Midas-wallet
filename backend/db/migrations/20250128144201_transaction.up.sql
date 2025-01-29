CREATE TABLE transaction (
    tx_hash VARCHAR(66) PRIMARY KEY,
    from_address VARCHAR(42),
    to_address VARCHAR(42),
    amount DECIMAL(64, 0),
    timestamp TIMESTAMP
    status VARCHAR(20)
    fee NUMERIC(18, 9)
);

CREATE indx_from_address ON transaction(from_address);
CREATE indx_to_address ON transaction(to_address);