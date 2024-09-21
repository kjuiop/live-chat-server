package mysql

type memoryClient struct {
}

func NewMemoryClient() Client {
	return &memoryClient{}
}

func (m memoryClient) GetServerList(qs string) ([]map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (m memoryClient) Close() {
	//TODO implement me
	panic("implement me")
}
