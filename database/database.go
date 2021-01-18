package database

import (
	"context"
	"time"

	pg "github.com/go-pg/pg/v10"
)

var conn *pg.DB

//Connect ...
func Connect() (*pg.DB, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if conn != nil {
		err := conn.Ping(ctx)
		if err == nil {
			return conn, nil
		}
	}

	opt, err := pg.ParseURL("")
	if err != nil {
		return nil, err
	}

	conn = pg.Connect(opt)

	return conn, nil
}
