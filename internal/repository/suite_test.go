package repository_test

import (
	"bulletin-board-api/internal/repository"
	"bytes"
	"database/sql/driver"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	testErrorMessage    = "test sql error"
	testAddressLocation = "test_address_location"
)

type gormCallSetup struct {
	want        bool
	expectation func(mock sqlmock.Sqlmock)
}

type AnyPgArray struct{}

func (a AnyPgArray) Match(v driver.Value) bool {
	value, ok := v.(string)
	var b pq.StringArray
	err := b.Scan(value)

	return err == nil && ok
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnyUUID struct{}

func (a AnyUUID) Match(v driver.Value) bool {
	_, err := uuid.Parse(v.(string))
	return err == nil
}

type AddressRepoTester struct {
	mockGorm sqlmock.Sqlmock

	underTesting repository.AddressRepository
}

func NewAddressRepoTester(t *testing.T) *AddressRepoTester {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	postgresDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	assert.NoError(t, err)
	repo := repository.NewAddressRepository(postgresDB, zap.NewNop())

	return &AddressRepoTester{
		mockGorm:     mock,
		underTesting: repo,
	}
}

func (f *AddressRepoTester) AssertAll(t *testing.T) {
	t.Helper()
	err := f.mockGorm.ExpectationsWereMet()
	assert.NoError(t, err)
}

type VaspRepoTester struct {
	underTesting repository.VaspRepository
}

func NewVaspRepoTester() *VaspRepoTester {
	viper.SetConfigType("yaml")

	prodConfig, err := os.ReadFile("../../configs/vasp/production.yaml")

	if err != nil {
		panic(err)
	}

	if err := viper.ReadConfig(bytes.NewReader(prodConfig)); err != nil {
		panic(err)
	}
	return &VaspRepoTester{underTesting: repository.NewVaspRepository()}
}

func NewVaspRepoTesterPreprod() *VaspRepoTester {
	viper.SetConfigType("yaml")

	prodConfig, err := os.ReadFile("../../configs/vasp/preprod.yaml")

	if err != nil {
		panic(err)
	}

	if err := viper.ReadConfig(bytes.NewReader(prodConfig)); err != nil {
		panic(err)
	}
	return &VaspRepoTester{underTesting: repository.NewVaspRepository()}
}
