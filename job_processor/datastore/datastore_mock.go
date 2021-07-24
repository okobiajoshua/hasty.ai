package datastore

import "github.com/stretchr/testify/mock"

type MockDataStore struct {
	mock.Mock
}

func NewMockDataStore() *MockDataStore {
	return &MockDataStore{}
}

func (m *MockDataStore) Update(objectID, status string) error {
	args := m.Called(objectID, status)
	return args.Error(0)
}
