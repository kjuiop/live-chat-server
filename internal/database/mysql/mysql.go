package mysql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"live-chat-server/config"
	"live-chat-server/internal/domain/system"
	"live-chat-server/internal/utils"
	"log"
	"log/slog"
	"strings"
)

const sqlFilePath = "internal/database/mysql/schema/database.sql"

type mysqlClient struct {
	cfg config.Mysql
	db  *sql.DB
}

func NewMysqlSingleClient(ctx context.Context, cfg config.Mysql) (Client, error) {

	url := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Database)
	db, err := sql.Open(cfg.Driver, url)
	if err != nil {
		return nil, fmt.Errorf("failed connect mysql, err : %w", err)
	}

	client := &mysqlClient{
		cfg: cfg,
		db:  db,
	}

	if err := client.checkDefaultTable(); err != nil {
		return nil, err
	}

	return client, nil
}

func (m *mysqlClient) Close() {
	if err := m.db.Close(); err != nil {
		slog.Error("failed db close", "error", err.Error())
	}
}

func (m *mysqlClient) checkDefaultTable() error {

	query := checkExistChatQuery()

	var count int
	if err := m.db.QueryRow(query).Scan(&count); err != nil {
		log.Fatal("Error checking table existence: ", err)
	}

	if count > 0 {
		return nil
	}

	content, err := utils.ReadFileContent(sqlFilePath)
	if err != nil {
		log.Fatalf("error opening SQL file: %v", err)
	}

	// SQL 파일 내용을 문자열로 변환하고 쿼리를 세미콜론(;)으로 분리
	queries := strings.Split(content, ";")

	// 각 쿼리를 실행
	for _, query := range queries {
		// 공백 제거
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		_, err = m.db.Exec(query)
		if err != nil {
			return fmt.Errorf("error executing query : %s, err : %w", query, err)
		}
	}

	return nil
}

func (m *mysqlClient) GetServerList(qs string) ([]system.ServerInfo, error) {

	cursor, err := m.db.Query(qs)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var result []system.ServerInfo

	for cursor.Next() {
		d := new(system.ServerInfo)

		if err := cursor.Scan(
			d.IP,
			d.Available,
		); err != nil {
			return nil, err
		}

		result = append(result, *d)
	}

	if len(result) == 0 {
		return []system.ServerInfo{}, nil
	}

	return result, nil
}

func checkExistChatQuery() string {
	return `
    SELECT COUNT(*)
    FROM information_schema.tables
    WHERE table_schema = 'chatting' AND table_name = 'chat';
    `
}