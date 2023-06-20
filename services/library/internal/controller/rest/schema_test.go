package rest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/controller/rest"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/library"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/model"
)

type MockSchemaService struct {
	schema *model.Schema
	err    error
}

func (s *MockSchemaService) Get(_ string) (interface{}, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.schema, nil
}

func (s *MockSchemaService) Search() (*model.Schemas, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &model.Schemas{
		&model.Schema{
			Title:       "TestSchema1",
			Description: "This is a test schema 1",
			Name:        "TestSchema1",
			URL:         "https://test.com/TestSchema1",
		},
		&model.Schema{
			Title:       "TestSchema2",
			Description: "This is a test schema 2",
			Name:        "TestSchema2",
			URL:         "https://test.com/TestSchema2",
		},
	}, nil
}

func TestSchemaHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		schemaName     string
		mockSvc        *MockSchemaService
		expectedStatus int
	}{
		{
			name:       "schema not found",
			schemaName: "NonexistentSchema",
			mockSvc: &MockSchemaService{
				err: library.SchemaNotFoundError{
					SchemaName: "NonexistentSchema",
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:       "database error",
			schemaName: "SchemaCausingDBError",
			mockSvc: &MockSchemaService{
				err: library.DatabaseError{},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:       "unknown error",
			schemaName: "SchemaCausingUnknownError",
			mockSvc: &MockSchemaService{
				err: errors.New("unknown error"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:       "success",
			schemaName: "TestSchema",
			mockSvc: &MockSchemaService{
				schema: &model.Schema{
					Name:        "TestSchema",
					Description: "This is a test schema",
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := rest.NewSchemaHandler(tt.mockSvc)

			r := gin.Default()
			r.GET("/schemas/:schemaName", handler.Get)

			req, _ := http.NewRequest(
				http.MethodGet,
				"/schemas/"+tt.schemaName,
				nil,
			)
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			require.Equal(t, tt.expectedStatus, resp.Code)
		})
	}
}

func TestSchemaHandler_Search(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSvc        *MockSchemaService
		expectedStatus int
	}{
		{
			name:           "success",
			mockSvc:        &MockSchemaService{},
			expectedStatus: http.StatusOK,
		},
		{
			name: "database error",
			mockSvc: &MockSchemaService{
				err: &library.DatabaseError{},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "unknown error",
			mockSvc: &MockSchemaService{
				err: errors.New("unknown error"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := rest.NewSchemaHandler(tt.mockSvc)

			r := gin.Default()
			r.GET("/schemas", handler.Search)

			req, _ := http.NewRequest(http.MethodGet, "/schemas", nil)
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			require.Equal(t, tt.expectedStatus, resp.Code)
		})
	}
}
