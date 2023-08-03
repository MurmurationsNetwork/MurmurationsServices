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
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/model"
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
func (s *SchemaParser) convertToBson(data []byte) (bson.D, error) {
	fullData := orderedmap.New()

	err := json.Unmarshal(data, &fullData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	bsonData := s.parseProperties(*fullData)

	return bsonData, nil
}

// parseProperties recursively transforms an ordered map representing a schema
// into its corresponding BSON document.
func (s *SchemaParser) parseProperties(fullData orderedmap.OrderedMap) bson.D {
	properties, exist := fullData.Get(PropertyKey)
	if !exist {
		return s.convertToBsonDocument(fullData)
	}

	propertiesMap := properties.(orderedmap.OrderedMap)

	for _, key := range propertiesMap.Keys() {
		value, _ := propertiesMap.Get(key)
		valueMap := value.(orderedmap.OrderedMap)

		var bsonDoc bson.D

		refPath, hasRef := valueMap.Get(ReferenceKey)
		if hasRef && refPath != nil {
			subSchema, err := s.fetchReferencedSchema(refPath.(string))
			if err != nil {
				continue
			}
			bsonDoc = s.parseProperties(*subSchema)
			bsonDoc = s.applyOverrides(bsonDoc, valueMap)
		} else {
			refType, hasType := valueMap.Get(TypeKey)
			if hasType && refType == ArrayType {
				bsonDoc = s.convertArrayToBsonDocument(valueMap)
			} else {
				bsonDoc = s.convertToBsonDocument(valueMap)
			}
		}

		propertiesMap.Set(key, bsonDoc)
	}

	propertiesData := s.convertToBsonDocument(propertiesMap)

	fullData.Set(PropertyKey, propertiesData)

	return s.convertToBsonDocument(fullData)
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
) (*orderedmap.OrderedMap, error) {
	_, fieldName := path.Split(url)

	fieldListURL := s.FieldListMap[fieldName]

	if fieldListURL == "" {
		return nil, fmt.Errorf(
			"get schema failed, url: %s, fieldListURL is empty",
			url,
		)
	}

	fieldJSON, err := s.fetchSchema(fieldListURL)
	if err != nil {
		return nil, err
	}

	subSchema := orderedmap.New()
	err = json.Unmarshal(fieldJSON, &subSchema)
	if err != nil {
		return nil, err
	}

	return subSchema, nil
}

// convertArrayToBsonDocument parses an array-like ordered map to a BSON document.
func (s *SchemaParser) convertArrayToBsonDocument(
	orderedMap orderedmap.OrderedMap,
) bson.D {
	items, _ := orderedMap.Get(ItemsKey)
	itemsMap := items.(orderedmap.OrderedMap)
	parsedSubSchema := s.parseProperties(itemsMap)

	bsonDoc := bson.D{}
	for _, key := range orderedMap.Keys() {
		var bsonElement bson.E

		if key == ItemsKey {
			bsonElement = bson.E{Key: key, Value: parsedSubSchema}
		} else {
			value, _ := orderedMap.Get(key)
			bsonElement = bson.E{Key: key, Value: value}
		}

		bsonDoc = append(bsonDoc, bsonElement)
	}

	return bsonDoc
}

// convertToBsonDocument parses an ordered map into a BSON document.
func (s *SchemaParser) convertToBsonDocument(
	orderedMap orderedmap.OrderedMap,
) bson.D {
	bsonData := bson.D{}

	for _, key := range orderedMap.Keys() {
		value, _ := orderedMap.Get(key)

		orderedMap, isOrderedMap := value.(orderedmap.OrderedMap)
		if !isOrderedMap {
			bsonData = append(bsonData, bson.E{Key: key, Value: value})
		} else {
			bsonData = append(bsonData, bson.E{
				Key: key,
				// Recursive call to handle nested maps.
				Value: s.parseProperties(orderedMap),
			})
		}
	}

	return bsonData
}
