-- +goose NO TRANSACTION
-- +goose Up
CREATE INDEX CONCURRENTLY IF NOT EXISTS index_address_ownership_proofs_iou ON address_ownership_proofs (chain, iou) WHERE iou = true AND deleted_at IS NULL;

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS index_address_ownership_proofs_iou;


