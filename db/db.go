// db/db.go
package db

import (
	"context"
	"github.com/jackc/pgx/v4"
)

// ConnectDB - 데이터베이스 연결 설정 함수
func ConnectDB(dbURL string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
