-- +goose Up
-- +goose StatementBegin
ALTER TABLE address_ownership_proofs ADD COLUMN encrypted_vasp_uuid VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE address_ownership_proofs DROP COLUMN encrypted_vasp_uuid;
-- +goose StatementEnd
