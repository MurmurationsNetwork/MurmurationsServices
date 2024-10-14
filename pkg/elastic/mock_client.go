package elastic

import (
	elastic "github.com/olivere/elastic/v7"
)

type mockClient struct {
}

func (*mockClient) setClient(_ *elastic.Client) {
}

func (*mockClient) GetClient() *elastic.Client {
	return nil
}

func (mockClient) Ping() error {
	return nil
}

func (*mockClient) CreateMappings(_ []Index) error {
	return nil
}

func (*mockClient) Index(
	_ string,
	_ interface{},
) (*elastic.IndexResponse, error) {
	return nil, nil
}

func (*mockClient) IndexWithID(
	_ string,
	_ string,
	_ interface{},
) (*elastic.IndexResponse, error) {
	return nil, nil
}

func (*mockClient) Search(
	_ string,
	_ *Query,
) (*elastic.SearchResult, error) {
	return nil, nil
}

func (*mockClient) Update(
	_ string,
	_ string,
	_ map[string]interface{},
) error {
	return nil
}

func (*mockClient) UpdateMany(
	_ string,
	_ *Query,
	_ map[string]interface{},
) error {
	return nil
}

func (*mockClient) Delete(_ string, _ string) error {
	return nil
}

func (*mockClient) DeleteMany(_ string, _ *Query) error {
	return nil
}

func (*mockClient) Export(
	_ string,
	_ *Query,
	_ []interface{},
) (*elastic.SearchResult, error) {
	return nil, nil
}

func (*mockClient) GetNodes(
	_ string,
	_ *Query,
) (*elastic.SearchResult, error) {
	return nil, nil
}
