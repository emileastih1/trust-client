package repository_test

import (
	"context"
	"encoding/pem"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestVaspRepository_GetProduction(t *testing.T) {
	t.Run("get vasp success", func(t *testing.T) {
		te := NewVaspRepoTester()
		uid, _ := uuid.Parse("bede7ca4-8a96-4f7f-876b-995b33212c83")
		vasp := te.underTesting.Get(context.Background(), uid)
		assert.NotEmpty(t, vasp)
	})

	t.Run("get vasp not found", func(t *testing.T) {
		te := NewVaspRepoTester()
		uid, _ := uuid.Parse("92837e40-772a-4c38-bece-39a16655ad95")
		vasp := te.underTesting.Get(context.Background(), uid)
		assert.Nil(t, vasp)
	})

	t.Run("check required fields", func(t *testing.T) {
		te := NewVaspRepoTester()
		for _, vasp := range te.underTesting.GetAll(context.Background()) {
			assert.NotEmpty(t, vasp.Name)
			assert.NotEmpty(t, vasp.Domain)
			assert.NotEmpty(t, vasp.PublicKey)
			assert.NotEmpty(t, vasp.PublicKeyVersion)
		}
	})

	t.Run("check public key is valid", func(t *testing.T) {
		te := NewVaspRepoTester()
		for _, vasp := range te.underTesting.GetAll(context.Background()) {
			block, _ := pem.Decode([]byte(vasp.PublicKey))
			assert.NotNil(t, block)
		}
	})

	t.Run("check https prefix is present in PII endpoint", func(t *testing.T) {
		te := NewVaspRepoTesterPreprod()
		for _, vasp := range te.underTesting.GetAll(context.Background()) {
			assert.Contains(t, vasp.PIIEndpoint, "https://")
		}
	})

	t.Run("check https prefix is present in PII request endpoint", func(t *testing.T) {
		te := NewVaspRepoTesterPreprod()
		for _, vasp := range te.underTesting.GetAll(context.Background()) {
			if vasp.PIIRequestEndpoint == "" {
				continue
			}
			assert.Contains(t, vasp.PIIRequestEndpoint, "https://")
		}
	})
}

func TestVaspRepository_GetPreprod(t *testing.T) {
	t.Run("get vasp success", func(t *testing.T) {
		te := NewVaspRepoTesterPreprod()
		uid, _ := uuid.Parse("6ac80f1b-cdf6-4e8f-99b3-ab12e5560f19")
		vasp := te.underTesting.Get(context.Background(), uid)
		assert.NotEmpty(t, vasp)
	})

	t.Run("get vasp not found", func(t *testing.T) {
		te := NewVaspRepoTesterPreprod()
		uid, _ := uuid.Parse("92837e40-772a-4c38-bece-39a16655ad95")
		vasp := te.underTesting.Get(context.Background(), uid)
		assert.Nil(t, vasp)
	})

	t.Run("check required fields", func(t *testing.T) {
		te := NewVaspRepoTesterPreprod()
		for _, vasp := range te.underTesting.GetAll(context.Background()) {
			assert.NotEmpty(t, vasp.Name)
			assert.NotEmpty(t, vasp.Domain)
			assert.NotEmpty(t, vasp.PublicKey)
			assert.NotEmpty(t, vasp.PublicKeyVersion)
		}
	})

	t.Run("check public key is valid", func(t *testing.T) {
		te := NewVaspRepoTesterPreprod()
		for _, vasp := range te.underTesting.GetAll(context.Background()) {
			if vasp.PublicKey == "PREPROD_DUMMY_KEY" {
				continue
			}
			block, _ := pem.Decode([]byte(vasp.PublicKey))
			if !assert.NotNil(t, block) {
				t.Logf("Public key is not valid for Vasp: %s", vasp.Name)
			}
		}
	})

	t.Run("check https prefix is present in PII endpoint", func(t *testing.T) {
		te := NewVaspRepoTesterPreprod()
		for _, vasp := range te.underTesting.GetAll(context.Background()) {
			assert.Contains(t, vasp.PIIEndpoint, "https://")
		}
	})

	t.Run("check https prefix is present in PII request endpoint", func(t *testing.T) {
		te := NewVaspRepoTesterPreprod()
		for _, vasp := range te.underTesting.GetAll(context.Background()) {
			if vasp.PIIRequestEndpoint == "" {
				continue
			}
			assert.Contains(t, vasp.PIIRequestEndpoint, "https://")
		}
	})

}
