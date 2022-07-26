/* ---- TOKEN PRICE ---- */
CREATE TABLE token_price
(
    id         SERIAL                      NOT NULL PRIMARY KEY,
    unit_name  TEXT                        NOT NULL REFERENCES token_unit (denom) UNIQUE,
    price      DECIMAL                     NOT NULL,
    market_cap BIGINT                      NOT NULL,
    timestamp  TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX token_price_timestamp_index ON token_price (timestamp);
