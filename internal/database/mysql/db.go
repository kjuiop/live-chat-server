package mysql

type Client interface {
	ExecQueryAndFetchRows(qs string, args ...interface{}) ([]map[string]interface{}, error)
	ExecQuery(qs string, args ...interface{}) error
	Close()
}
