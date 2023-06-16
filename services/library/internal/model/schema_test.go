package model_test

import (
	"testing"

	"github.com/iancoleman/orderedmap"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/model"
)

func TestSchema(t *testing.T) {
	schema := &model.Schema{
		Title:       "Title1",
		Description: "Description1",
		Name:        "Name1",
		URL:         "URL1",
	}
	want := &model.Schema{
		Title:       "Title1",
		Description: "Description1",
		Name:        "Name1",
		URL:         "URL1",
	}
	require.Equal(t, want, schema.Marshall())
}

func TestSchemas(t *testing.T) {
	schemas := model.Schemas{
		&model.Schema{
			Title:       "Title1",
			Description: "Description1",
			Name:        "Name1",
			URL:         "URL1",
		},
		&model.Schema{
			Title:       "Title2",
			Description: "Description2",
			Name:        "Name2",
			URL:         "URL2",
		},
	}
	want := []interface{}{
		&model.Schema{
			Title:       "Title1",
			Description: "Description1",
			Name:        "Name1",
			URL:         "URL1",
		},
		&model.Schema{
			Title:       "Title2",
			Description: "Description2",
			Name:        "Name2",
			URL:         "URL2",
		},
	}
	require.Equal(t, want, schemas.Marshall())
}

func TestSingleSchema(t *testing.T) {
	tests := []struct {
		name         string
		singleSchema *model.SingleSchema
		expected     *orderedmap.OrderedMap
	}{
		{
			name:         "Empty schema",
			singleSchema: &model.SingleSchema{},
			expected:     orderedmap.New(),
		},
		{
			name: "Schema with one simple element",
			singleSchema: &model.SingleSchema{
				FullSchema: bson.D{
					{Key: "key", Value: "value"},
				},
			},
			expected: func() *orderedmap.OrderedMap {
				m := orderedmap.New()
				m.Set("key", "value")
				return m
			}(),
		},
		{
			name: "Schema with nested element",
			singleSchema: &model.SingleSchema{
				FullSchema: bson.D{
					{Key: "key", Value: "value"},
					{Key: "nested", Value: bson.D{
						{Key: "nestedKey", Value: "nestedValue"},
					}},
				},
			},
			expected: func() *orderedmap.OrderedMap {
				m := orderedmap.New()
				m.Set("key", "value")
				nested := orderedmap.New()
				nested.Set("nestedKey", "nestedValue")
				m.Set("nested", nested)
				return m
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.singleSchema.ToMap()
			require.Equal(t, tt.expected, result)
		})
	}
}
