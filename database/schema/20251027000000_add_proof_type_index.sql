-- +goose NO TRANSACTION
-- +goose Up
CREATE INDEX CONCURRENTLY IF NOT EXISTS index_address_ownership_proofs_proof_type ON address_ownership_proofs (chain, proof_type) WHERE proof_type IS NOT NULL AND deleted_at IS NULL;

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS index_address_ownership_proofs_proof_type;

