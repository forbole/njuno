/* ---- BLOCK ---- */
CREATE TABLE block
(
    height           BIGINT  UNIQUE PRIMARY KEY,
    hash             TEXT    NOT NULL UNIQUE,
    num_txs          INTEGER DEFAULT 0,
    total_gas        BIGINT  DEFAULT 0,
    proposer_address TEXT    NOT NULL,
    timestamp        TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX block_height_index ON block (height);
CREATE INDEX block_hash_index ON block (hash);
CREATE INDEX block_proposer_address_index ON block (proposer_address);


/* ---- PRE COMMIT ---- */
CREATE TABLE pre_commit
(
    validator_address TEXT                        NOT NULL,
    height            BIGINT                      NOT NULL REFERENCES block (height),
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    voting_power      BIGINT                      NOT NULL,
    proposer_priority BIGINT                      NOT NULL,
    UNIQUE (validator_address, timestamp)
);
CREATE INDEX pre_commit_validator_address_index ON pre_commit (validator_address);
CREATE INDEX pre_commit_height_index ON pre_commit (height);


/* ---- TRANSACTION ---- */
CREATE TABLE transaction
(
    hash         TEXT    NOT NULL,
    height       BIGINT  NOT NULL REFERENCES block (height),

    /* Body */
    memo         TEXT,
    signatures   TEXT[],
    fee          COIN[] NOT NULL DEFAULT '{}',
    gas          TEXT,

    /* PSQL partition */
    partition_id BIGINT  NOT NULL DEFAULT 0,

    CONSTRAINT unique_tx UNIQUE (hash, partition_id)
) PARTITION BY LIST (partition_id);
CREATE INDEX transaction_hash_index ON transaction (hash);
CREATE INDEX transaction_height_index ON transaction (height);
CREATE INDEX transaction_partition_id_index ON transaction (partition_id);