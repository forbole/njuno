/* ---- VALIDATOR INFO ---- */
CREATE TABLE validator
(
    consensus_address     TEXT   NOT NULL UNIQUE,
    self_delegate_address TEXT   NOT NULL PRIMARY KEY,
    height                BIGINT NOT NULL
);
CREATE INDEX validator_self_delegate_address_index ON validator (self_delegate_address);


/* ---- VALIDATOR VOTING POWER ---- */
CREATE TABLE validator_voting_power
(
    validator_address     TEXT   NOT NULL UNIQUE,
    self_delegate_address TEXT   NOT NULL PRIMARY KEY REFERENCES validator (self_delegate_address),
    voting_power          TEXT   NOT NULL DEFAULT '0',
    height                BIGINT NOT NULL
);
CREATE INDEX validator_voting_power_height_index ON validator_voting_power (height);


/* ---- VALIDATOR DESCRIPTION ---- */
CREATE TABLE validator_description
(
    validator_address     TEXT   NOT NULL UNIQUE,
    self_delegate_address TEXT   NOT NULL PRIMARY KEY REFERENCES validator (self_delegate_address),
    moniker               TEXT,
    identity              TEXT,
    avatar_url            TEXT,
    details               TEXT,
    height                BIGINT NOT NULL
);
CREATE INDEX validator_description_height_index ON validator_description (height);


/* ---- VALIDATOR COMMISSION ---- */
CREATE TABLE validator_commission
(
    validator_address     TEXT   NOT NULL UNIQUE,
    self_delegate_address TEXT   NOT NULL PRIMARY KEY REFERENCES validator (self_delegate_address),
    commission            TEXT   NOT NULL,
    min_self_delegation   TEXT   NOT NULL DEFAULT '0',
    height                BIGINT NOT NULL
);
CREATE INDEX validator_commission_height_index ON validator_commission (height);

/* ---- VALIDATOR STATUS ---- */
CREATE TABLE validator_status
(
    validator_address     TEXT   NOT NULL UNIQUE,
    self_delegate_address TEXT   NOT NULL PRIMARY KEY REFERENCES validator (self_delegate_address),
    in_active_set         TEXT   NOT NULL,
    jailed                TEXT   NOT NULL,
    tombstoned            TEXT   NOT NULL,
    height                BIGINT NOT NULL
);
CREATE INDEX validator_status_height_index ON validator_status (height);

/* ---- DOUBLE SIGN VOTE ---- */
CREATE TABLE double_sign_vote
(
    id                SERIAL PRIMARY KEY,
    type              SMALLINT NOT NULL,
    height            BIGINT   NOT NULL,
    round             INT      NOT NULL,
    block_id          TEXT     NOT NULL,
    validator_address TEXT     NOT NULL REFERENCES validator (consensus_address),
    validator_index   INT      NOT NULL,
    signature         TEXT     NOT NULL,
    UNIQUE (block_id, validator_address)
);
CREATE INDEX double_sign_vote_validator_address_index ON double_sign_vote (validator_address);
CREATE INDEX double_sign_vote_height_index ON double_sign_vote (height);


/* ---- DOUBLE SIGN EVIDENCE ---- */
CREATE TABLE double_sign_evidence
(
    height    BIGINT NOT NULL,
    vote_a_id BIGINT NOT NULL REFERENCES double_sign_vote (id),
    vote_b_id BIGINT NOT NULL REFERENCES double_sign_vote (id)
);
CREATE INDEX double_sign_evidence_height_index ON double_sign_evidence (height);
