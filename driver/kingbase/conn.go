package kingbase

import (
	"context"
	"database/sql/driver"
)

type KBConn struct {
	conn driver.Conn
}

func (KC KBConn) ExecContext(ctx context.Context, sql string, args []driver.NamedValue) (driver.Result, error) {
	return KC.conn.(driver.ExecerContext).ExecContext(ctx, sql, args)
}

func (KC KBConn) QueryContext(ctx context.Context, sql string, args []driver.NamedValue) (driver.Rows, error) {
	return KC.conn.(driver.QueryerContext).QueryContext(ctx, sql, args)
}

func (KC KBConn) PrepareContext(ctx context.Context, sql string) (driver.Stmt, error) {
	return KC.conn.Prepare(sql)
}

func (KC KBConn) Prepare(sql string) (driver.Stmt, error) {
	return KC.conn.Prepare(sql)
}

func (KC KBConn) Begin() (driver.Tx, error) {
	return KC.conn.Begin()
}

func (KC KBConn) Close() error {
	return KC.conn.Close()
}
