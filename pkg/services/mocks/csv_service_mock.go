// pkg/services/mocks/csv_service_mock.go
package mocks

import (
	"encoding/csv"

	"github.com/stretchr/testify/mock"
)

type MockCSVService struct {
	mock.Mock
}

func (m *MockCSVService) ReadCSVHeaders(reader *csv.Reader) ([]string, error) {
	args := m.Called(reader)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCSVService) ReadCSVRecord(reader *csv.Reader) ([]string, error) {
	args := m.Called(reader)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCSVService) GetCSVHeaders(model interface{}) []string {
	args := m.Called(model)
	return args.Get(0).([]string)
}
