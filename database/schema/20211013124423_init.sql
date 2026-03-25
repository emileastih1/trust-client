-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS address_ownership_proofs
(
    id                   UUID NOT NULL PRIMARY KEY,
    address              TEXT NOT NULL,
    chain                VARCHAR(128) NOT NULL,
    signature            TEXT,
    iou                  BOOLEAN DEFAULT NULL,
    created_at           TIMESTAMP NOT NULL,
    updated_at           TIMESTAMP NOT NULL,
    deleted_at           TIMESTAMP DEFAULT NULL,
    prefix               TEXT,
    proof_submitted_time TIMESTAMP,
    proof_type           VARCHAR(256),
    aux_proof_data       JSON
);

CREATE TABLE IF NOT EXISTS cuckoo_filter_bucket
(
    id         BIGSERIAL NOT NULL PRIMARY KEY,
    reference  BIGINT NOT NULL,
    vasp_uuid  UUID NOT NULL,
    filled     INTEGER DEFAULT 0,
    slots      TEXT[],
    chain      TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_index_address_ownership_proofs_active ON address_ownership_proofs (chain, address, (deleted_at IS NULL)) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uniq_index_cuckoo_filter_bucket ON cuckoo_filter_bucket (chain, vasp_uuid, reference, (deleted_at IS NULL)) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS address_ownership_proofs;
DROP TABLE IF EXISTS cuckoo_filter_bucket;
-- +goose StatementEnd
