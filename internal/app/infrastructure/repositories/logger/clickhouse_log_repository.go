package logger

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseRepository struct {
	conn  clickhouse.Conn
	table string
}

func NewClickHouseRepository(dsn, table string) (*ClickHouseRepository, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{
			dsn,
		},
	})
	if err != nil {
		return nil, err
	}

	return &ClickHouseRepository{
		conn:  conn,
		table: table,
	}, nil
}

func (c *ClickHouseRepository) Save(ctx context.Context, log map[string]any) error {
	return c.conn.Exec(
		ctx,
		"INSERT INTO "+c.table+" (timestamp, level, msg, fields) VALUES (?, ?, ?, ?)",
		log["timestamp"],
		log["level"],
		log["msg"],
		log["fields"],
	)
}
