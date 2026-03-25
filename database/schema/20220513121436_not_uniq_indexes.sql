-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS index_address_ownership_proofs ON address_ownership_proofs (chain, address, deleted_at);
CREATE INDEX IF NOT EXISTS index_cuckoo_filter_bucket ON cuckoo_filter_bucket (chain, reference, vasp_uuid, deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS index_address_ownership_proofs;
DROP INDEX IF EXISTS index_cuckoo_filter_bucket;
-- +goose StatementEnd
