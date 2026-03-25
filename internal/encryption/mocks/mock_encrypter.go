package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockEncrypter struct {
	mock.Mock
}

func (m *MockEncrypter) Encrypt(ctx context.Context, value string) (string, error) {
	args := m.Called(ctx, value)
	return args.String(0), args.Error(1)
}

func (m *MockEncrypter) Decrypt(ctx context.Context, value string) (string, error) {
	args := m.Called(ctx, value)
	return args.String(0), args.Error(1)
}
