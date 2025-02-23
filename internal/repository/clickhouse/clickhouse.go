package clickhouse

import (
	"context"
	"log"
	"net"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"gymnote/internal/config"
)

type clickHouse struct {
	conn driver.Conn
	cfg  *config.DBConfig
}

func New(ctx context.Context, cfg *config.DBConfig) (*clickHouse, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{net.JoinHostPort(cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Name,
			Username: cfg.User,
			Password: cfg.Password,
		},
		Debugf: func(format string, v ...any) {
			log.Printf(format, v)
		},
	})
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}

	return &clickHouse{conn: conn, cfg: cfg}, nil
}

func (c *clickHouse) Close() {
	c.conn.Close()
}
