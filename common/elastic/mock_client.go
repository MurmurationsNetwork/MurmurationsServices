package elastic

type mockClient struct {
}

func (c *mockClient) setClient(client *elastic.Client) {
}

func (c *mockClient) CreateMappings(indices []Index) error {
	return nil
}

func (c *mockClient) Index(index string, doc interface{}) (*elastic.IndexResponse, error) {
	return nil, nil
}

func (c *mockClient) IndexWithID(index string, id string, doc interface{}) (*elastic.IndexResponse, error) {
	return nil, nil
}

func (c *mockClient) Search(index string, query *Query) (*elastic.SearchResult, error) {
	return nil, nil
}

func (c *mockClient) Delete(index string, id string) error {
	return nil
}
