/* ---- GENESIS ---- */
CREATE TABLE genesis
(
    one_row_id     BOOL      NOT NULL DEFAULT TRUE PRIMARY KEY,
    chain_id       TEXT      NOT NULL,
    time           TIMESTAMP NOT NULL,
    initial_height BIGINT    NOT NULL,
    CHECK (one_row_id)
);


/* ---- AVERAGE BLOCK PER MINUTE ---- */
CREATE TABLE average_block_time_per_minute
(
    one_row_id   BOOL    NOT NULL DEFAULT TRUE PRIMARY KEY,
    average_time DECIMAL NOT NULL,
    height       BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX average_block_time_per_minute_height_index ON average_block_time_per_minute (height);


/* ----  AVERAGE BLOCK PER HOUR ---- */
CREATE TABLE average_block_time_per_hour
(
    one_row_id   BOOL    NOT NULL DEFAULT TRUE PRIMARY KEY,
    average_time DECIMAL NOT NULL,
    height       BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX average_block_time_per_hour_height_index ON average_block_time_per_hour (height);


/* ----  AVERAGE BLOCK PER DAY ---- */
CREATE TABLE average_block_time_per_day
(
    one_row_id   BOOL    NOT NULL DEFAULT TRUE PRIMARY KEY,
    average_time DECIMAL NOT NULL,
    height       BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX average_block_time_per_day_height_index ON average_block_time_per_day (height);


/* ----  AVERAGE BLOCK FROM GENESIS ---- */
CREATE TABLE average_block_time_from_genesis
(
    one_row_id   BOOL    NOT NULL DEFAULT TRUE PRIMARY KEY,
    average_time DECIMAL NOT NULL,
    height       BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX average_block_time_from_genesis_height_index ON average_block_time_from_genesis (height);
