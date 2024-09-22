package mysql

type Client interface {
	GetServerList(qs string) ([]map[string]interface{}, error)
	ExecQuery(qs string, args ...interface{}) error
	Close()
}
