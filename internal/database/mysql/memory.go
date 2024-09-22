package mysql

type memoryClient struct {
}

func NewMemoryClient() Client {
	return &memoryClient{}
}

func (m memoryClient) ExecQueryAndFetchRows(qs string, args ...interface{}) ([]map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (m memoryClient) ExecQuery(qs string, args ...interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m memoryClient) Close() {
	//TODO implement me
	panic("implement me")
}
