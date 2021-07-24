package datastore

import "github.com/stretchr/testify/mock"

type MockDataStore struct {
	mock.Mock
}

func NewMockDataStore() *MockDataStore {
	return &MockDataStore{}
}

func (m *MockDataStore) Save() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDataStore) GetByObjectID(objectID string) (interface{}, error) {
	args := m.Called(objectID)
	return nil, args.Error(1)
}

func (m *MockDataStore) GetByJobID(jobID string) (interface{}, error) {
	args := m.Called(jobID)
	return nil, args.Error(1)
}
