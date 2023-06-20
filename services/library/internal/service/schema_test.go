package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/service"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Get(schemaName string) (interface{}, error) {
	args := m.Called(schemaName)
	return args.Get(0), args.Error(1)
}

func (m *MockRepo) Search() (*model.Schemas, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Schemas), args.Error(1)
}

func TestSchemaService(t *testing.T) {
	mockSchema := &model.Schema{
		Title:       "Title",
		Description: "Description",
		Name:        "TestSchema",
		URL:         "URL",
	}
	mockSchemas := &model.Schemas{
		&model.Schema{
			Title:       "Title1",
			Description: "Description1",
			Name:        "TestSchema1",
			URL:         "URL1",
		},
		&model.Schema{
			Title:       "Title2",
			Description: "Description2",
			Name:        "TestSchema2",
			URL:         "URL2",
		},
	}

	tests := []struct {
		name            string
		repoGet         func(m *MockRepo)
		repoSearch      func(m *MockRepo)
		getSchemaName   string
		expGetErr       bool
		expSearchErr    bool
		expGetResult    interface{}
		expSearchResult *model.Schemas
	}{
		{
			name: "Test valid Get and Search",
			repoGet: func(m *MockRepo) {
				m.On("Get", "TestSchema").Return(mockSchema, nil)
			},
			repoSearch: func(m *MockRepo) {
				m.On("Search").Return(mockSchemas, nil)
			},
			getSchemaName:   "TestSchema",
			expGetErr:       false,
			expSearchErr:    false,
			expGetResult:    mockSchema,
			expSearchResult: mockSchemas,
		},
		{
			name: "Test Get and Search errors",
			repoGet: func(m *MockRepo) {
				m.On("Get", "NonExistent").
					Return(nil, errors.New("schema not found"))
			},
			repoSearch: func(m *MockRepo) {
				m.On("Search").Return(nil, errors.New("database error"))
			},
			getSchemaName:   "NonExistent",
			expGetErr:       true,
			expSearchErr:    true,
			expGetResult:    nil,
			expSearchResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			tt.repoGet(mockRepo)
			tt.repoSearch(mockRepo)
			s := service.NewSchemaService(mockRepo)

			resultGet, err := s.Get(tt.getSchemaName)
			if tt.expGetErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expGetResult, resultGet)
			}

			resultSearch, err := s.Search()
			if tt.expSearchErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expSearchResult, resultSearch)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
