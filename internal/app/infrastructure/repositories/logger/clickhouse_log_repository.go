package logger

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
)

const (
	TimestampKey     = "timestamp"
	LevelKey         = "level"
	MsgKey           = "msg"
	FieldsKey        = "fields"
	CorrelationIDKey = "correlation_id"
)

type ClickHouseRepository struct {
	conn  clickhouse.Conn
	table string
}

func NewClickHouseRepository(conn clickhouse.Conn, table string) *ClickHouseRepository {
	return &ClickHouseRepository{
		conn:  conn,
		table: table,
	}
}

func (c *ClickHouseRepository) Save(ctx context.Context, log map[string]any) error {
	query := "INSERT INTO " + c.table + " (timestamp, level, msg, fields, correlation_id) VALUES (?, ?, ?, ?, ?)"
	return c.conn.Exec(ctx, query,
		log[TimestampKey],
		log[LevelKey],
		log[MsgKey],
		log[FieldsKey],
		log[CorrelationIDKey],
	)
}
