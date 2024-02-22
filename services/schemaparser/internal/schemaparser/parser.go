package schemaparser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/iancoleman/orderedmap"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/schemaparser/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/schemaparser/internal/model"
)

const (
	// PropertyKey is used as the key to access properties in our ordered map.
	PropertyKey = "properties"
	// ReferenceKey is used as the key to fetch a reference to another schema in our ordered map.
	ReferenceKey = "$ref"
	// TypeKey is used as the key to determine the type of a value in our ordered map.
	TypeKey = "type"
	// ItemsKey is used to determine if a value is an items type.
	ItemsKey = "items"
	// ArrayType is used to determine if a value is an array type.
	ArrayType = "array"
	// ObjectType is used to determine if a value is an object type.
	ObjectType = "object"
)

// SchemaResult represents the get schema result.
type SchemaResult struct {
	// Parsed JSON schema data.
	Schema *model.SchemaJSON
	// Full schema data as BSON.
	FullJSON bson.D
}

// SchemaParser represents the schema parser.
type SchemaParser struct {
	// Field map list.
	FieldListMap map[string]string
}

// NewSchemaParser creates a new instance of SchemaParser.
func NewSchemaParser(fieldListMap map[string]string) *SchemaParser {
	return &SchemaParser{
		FieldListMap: fieldListMap,
	}
}

// GetSchema fetches, parses and converts a schema to BSON from a given URL.
func (s *SchemaParser) GetSchema(url string) (*SchemaResult, error) {
	schemaData, err := s.fetchSchema(url)
	if err != nil {
		return nil, err
	}

	parsedSchema, err := s.parseSchema(schemaData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	fullJSON, err := s.convertToBson(schemaData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert schema to BSON: %w", err)
	}

	return &SchemaResult{Schema: parsedSchema, FullJSON: fullJSON}, nil
}

// GetLocalSchema fetches, parses and converts a schema to BSON from local schema byte.
func (s *SchemaParser) GetLocalSchema(schema []byte, fields map[string][]byte) (*SchemaResult, error) {
	parsedSchema, err := s.parseSchema(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	fullJSON, err := s.convertToBson(schema, fields)
	if err != nil {
		return nil, fmt.Errorf("failed to convert schema to BSON: %w", err)
	}

	return &SchemaResult{Schema: parsedSchema, FullJSON: fullJSON}, nil
}

// fetchSchema fetches the schema data from the provided URL.
func (s *SchemaParser) fetchSchema(url string) ([]byte, error) {
	// Perform a GET request with bearer token authentication.
	resp, err := httputil.GetWithBearerToken(url, config.Values.Github.TOKEN)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get file failed, url: %s, status code: %d",
			url, resp.StatusCode)
	}

	// Read the entire body of the response.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the body into a map.
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	// Decode the base64-encoded content.
	decodedContent, err := base64.StdEncoding.DecodeString(
		data["content"].(string),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 content: %w", err)
	}

	// Check if the decoded content is empty or not.
	if len(decodedContent) == 0 {
		return nil, fmt.Errorf(
			"get file failed, url: %s, content is empty",
			url,
		)
	}

	return decodedContent, nil
}

// parseSchema parses the schema data into the SchemaJSON type.
func (s *SchemaParser) parseSchema(data []byte) (*model.SchemaJSON, error) {
	var schema model.SchemaJSON

	err := json.Unmarshal(data, &schema)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}

// convertToBson converts the schema data into a bson.D type.
func (s *SchemaParser) convertToBson(data []byte, optionalFields ...map[string][]byte) (bson.D, error) {
	fullData := orderedmap.New()

	err := json.Unmarshal(data, &fullData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	var bsonData bson.D
	if len(optionalFields) > 0 && optionalFields[0] != nil {
		bsonData, err = s.parseProperties(*fullData, optionalFields[0])
	} else {
		bsonData, err = s.parseProperties(*fullData)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse properties: %w", err)
	}

	return bsonData, nil
}

// parseProperties recursively transforms a map representing a schema into a BSON document.
func (s *SchemaParser) parseProperties(
	fullData orderedmap.OrderedMap,
	optionalFields ...map[string][]byte,
) (bson.D, error) {
	properties, exist := fullData.Get(PropertyKey)
	if !exist {
		return s.convertToBsonDocument(fullData)
	}

	propertiesMap, ok := properties.(orderedmap.OrderedMap)
	if !ok {
		return nil, fmt.Errorf(
			"unexpected type for properties. Expected map, got: %T",
			properties,
		)
	}

	for _, key := range propertiesMap.Keys() {
		value, _ := propertiesMap.Get(key)
		valueMap, ok := value.(orderedmap.OrderedMap)
		if !ok {
			return nil, fmt.Errorf(
				"unexpected type for value with key '%s'. Expected map, got: %T",
				key,
				value,
			)
		}

		var bsonDoc bson.D
		var err error

		refPath, hasRef := valueMap.Get(ReferenceKey)
		if hasRef && refPath != nil {
			path, ok := refPath.(string)
			if !ok {
				return nil, fmt.Errorf(
					"unexpected type for refPath. Expected string, got: %T",
					refPath,
				)
			}
			var subSchema *orderedmap.OrderedMap
			if len(optionalFields) > 0 && optionalFields[0] != nil {
				subSchema, err = s.fetchReferencedSchema(path, optionalFields[0])
			} else {
				subSchema, err = s.fetchReferencedSchema(path)
			}
			if err != nil {
				return nil, fmt.Errorf(
					"failed to fetch referenced schema with refPath '%s' and key '%s': %w",
					path,
					key,
					err,
				)
			}
			if len(optionalFields) > 0 && optionalFields[0] != nil {
				bsonDoc, err = s.parseProperties(*subSchema, optionalFields[0])
			} else {
				bsonDoc, err = s.parseProperties(*subSchema)
			}
			if err != nil {
				return nil, err
			}
			bsonDoc = s.applyOverrides(bsonDoc, valueMap)
		} else {
			refType, hasType := valueMap.Get(TypeKey)
			if hasType && refType == ArrayType {
				bsonDoc, err = s.convertArrayPropToBson(valueMap)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to convert array to BSON document with key '%s': %w",
						key, err,
					)
				}
			} else {
				bsonDoc, err = s.convertToBsonDocument(valueMap)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to convert to BSON document for key '%s': %w",
						key, err,
					)
				}
			}
		}
		propertiesMap.Set(key, bsonDoc)
	}

	propertiesData, err := s.convertToBsonDocument(propertiesMap)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to convert properties map to BSON document: %w",
			err,
		)
	}
	fullData.Set(PropertyKey, propertiesData)

	result, err := s.convertToBsonDocument(fullData)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to convert full data to BSON document: %w",
			err,
		)
	}
	return result, nil
}

func (s *SchemaParser) applyOverrides(
	bsonDoc bson.D,
	orderedMap orderedmap.OrderedMap,
) bson.D {
	// Iterate over the original properties to apply overrides.
	for _, overrideKey := range orderedMap.Keys() {
		if overrideKey != ReferenceKey {
			overrideValue, _ := orderedMap.Get(overrideKey)
			for i, elem := range bsonDoc {
				if elem.Key == overrideKey {
					bsonDoc[i].Value = overrideValue
					break
				}
			}
		}
	}
	return bsonDoc
}

func (s *SchemaParser) fetchReferencedSchema(
	url string,
	optionalFields ...map[string][]byte,
) (*orderedmap.OrderedMap, error) {
	_, fieldName := path.Split(url)

	var (
		fieldJSON []byte
		ok        bool
		err       error
	)
	if len(optionalFields) > 0 && optionalFields[0] != nil {
		fieldJSON, ok = optionalFields[0][fieldName]
		if !ok {
			return nil, fmt.Errorf(
				"get schema failed, url: %s, fieldListURL is empty",
				url,
			)
		}
	} else {
		fieldListURL := s.FieldListMap[fieldName]

		if fieldListURL == "" {
			return nil, fmt.Errorf(
				"get schema failed, url: %s, fieldListURL is empty",
				url,
			)
		}

		fieldJSON, err = s.fetchSchema(fieldListURL)
		if err != nil {
			return nil, err
		}
	}

	subSchema := orderedmap.New()
	err = json.Unmarshal(fieldJSON, &subSchema)
	if err != nil {
		return nil, err
	}

	return subSchema, nil
}

// convertArrayPropsToBson transforms an array-like property into a BSON document.
func (s *SchemaParser) convertArrayPropToBson(
	orderedMap orderedmap.OrderedMap,
) (bson.D, error) {
	// Retrieve the 'items' field from the provided ordered map.
	items, exists := orderedMap.Get(ItemsKey)
	if !exists {
		return nil, fmt.Errorf("missing '%s' key in the ordered map", ItemsKey)
	}

	// Attempt to type assert the 'items' field to an ordered map.
	itemsMap, ok := items.(orderedmap.OrderedMap)
	if !ok {
		return nil, fmt.Errorf(
			"unexpected type for '%s'. Expected map",
			ItemsKey,
		)
	}

	// Parse the sub-schema found in the 'items' field.
	parsedSubSchema, err := s.parseProperties(itemsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to parse properties: %w", err)
	}

	bsonDoc := bson.D{}
	// Iterate over the keys in the orderedMap to construct the BSON document.
	for _, key := range orderedMap.Keys() {
		value, _ := orderedMap.Get(key)

		var bsonElement bson.E
		if key == ItemsKey {
			bsonElement = bson.E{Key: key, Value: parsedSubSchema}
		} else {
			bsonElement = bson.E{Key: key, Value: value}
		}

		bsonDoc = append(bsonDoc, bsonElement)
	}

	return bsonDoc, nil
}

// convertToBsonDocument parses an ordered map into a BSON document.
func (s *SchemaParser) convertToBsonDocument(
	orderedMap orderedmap.OrderedMap,
) (bson.D, error) {
	bsonData := bson.D{}

	for _, key := range orderedMap.Keys() {
		value, _ := orderedMap.Get(key)

		orderedMapValue, isOrderedMap := value.(orderedmap.OrderedMap)
		if !isOrderedMap {
			bsonData = append(bsonData, bson.E{Key: key, Value: value})
		} else {
			parsedValue, err := s.parseProperties(orderedMapValue)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to parse properties for key '%s': %w",
					key, err,
				)
			}
			bsonData = append(bsonData, bson.E{Key: key, Value: parsedValue})
		}
	}

	return bsonData, nil
}
