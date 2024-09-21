package mysql

type Client interface {
	GetServerList(qs string) ([]map[string]interface{}, error)
	Close()
}
