# Installation
Simple install the package to your [$GOPATH](https://github.com/golang/go/wiki/GOPATH "GOPATH") with the [go tool](https://golang.org/cmd/go/ "go command") from shell:
```bash
$ export GOPRIVATE="devops.aishu.cn"
$ go get -u github.com/kweaver-ai/proton-rds-sdk-go/sqlx
```
Make sure [Git is installed](https://git-scm.com/downloads) on your machine and in your system's `PATH`.

# Usage
Reference (https://pkg.go.dev/database/sql#DB)
```
type DB
	func Open(driverName, dataSourceName string) (*DB, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Ping() error
	PingContext(ctx context.Context) error
	SetConnMaxLifetime(d time.Duration)
	SetMasterMaxOpenConns(n int)
	SetBackupMaxOpenConns(n int)
```
# Example
Examples are available in example directory.
```go
import (
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

// ...

```