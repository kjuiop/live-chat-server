package mysql

type memoryClient struct {
}

func NewMemoryClient() Client {
	return &memoryClient{}
}
