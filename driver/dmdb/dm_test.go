package dmdb

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

type TestDBInfo struct {
	Host     string
	Port     int
	Username string
	Password string
}

func getTestDBInfo() TestDBInfo {
	user := os.Getenv("DM_TEST_USER")                    // 例如设置为 "test_user"
	password := os.Getenv("DM_TEST_PASSWORD")            // 例如设置为 "test_pwd_123"
	host := os.Getenv("DM_TEST_HOST")                    // 例如设置为 "localhost"
	port, err := strconv.Atoi(os.Getenv("DM_TEST_PORT")) // 例如设置为 "3306"
	if err != nil {
		log.Fatalf("DM_TEST_PORT is not a number: %v", err)
	}
	return TestDBInfo{
		Host:     host,
		Port:     port,
		Username: user,
		Password: password,
	}
}

func TestOpen(t *testing.T) {
	Convey("Test dmdb.Open\n", t, func() {
		info := getTestDBInfo()
		Convey("Open fail,no slash\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)SYSDBA?timeout=10s", info.Username, info.Password, info.Host, info.Port)
			_, err := Open(dsn)
			assert.Equal(t, err, errors.New("invalid DSN: missing the slash separating the database name"))
		})
		Convey("Open fail,missing symbol\n", func() {
			dsn := fmt.Sprintf("%s:%stcp%s:%d/SYSDBA?timeout=10s", info.Username, info.Password, info.Host, info.Port)
			_, err := Open(dsn)
			assert.Equal(t, err, errors.New("invalid DSN: missing '@' or '(' or ')' separating the necessary parts"))
		})
		Convey("Open fail,change fail,invalid time unit suffix\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout=10xxs", info.Username, info.Password, info.Host, info.Port)
			_, err := Open(dsn)
			assert.Equal(t, err, errors.New("time: unknown unit \"xxs\" in duration \"10xxs\""))
		})
		Convey("Open success,test invalid param continue\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout", info.Username, info.Password, info.Host, info.Port)
			_, err := Open(dsn)
			assert.Equal(t, err, nil)
		})
		Convey("Open success,case param two\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?readTimeout=10s&timeout=10s", info.Username, info.Password, info.Host, info.Port)
			_, err := Open(dsn)
			assert.Equal(t, err, nil)
		})
		Convey("Open success\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout=10s&autocommit=true&readTimeout=10s", info.Username, info.Password, info.Host, info.Port)
			_, err := Open(dsn)
			assert.Equal(t, err, nil)
		})
	})
}

func TestOpenConnector(t *testing.T) {
	Convey("Test dmdb.OpenConnector\n", t, func() {
		info := getTestDBInfo()
		Convey("OpenConnector fail,no slash\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)SYSDBA?timeout=10s", info.Username, info.Password, info.Host, info.Port)
			_, err := OpenConnector(dsn)
			assert.Equal(t, err, errors.New("invalid DSN: missing the slash separating the database name"))
		})
		Convey("OpenConnector fail,missing symbol\n", func() {
			dsn := fmt.Sprintf("%s:%stcp%s:%d/SYSDBA?timeout=10s", info.Username, info.Password, info.Host, info.Port)
			_, err := OpenConnector(dsn)
			assert.Equal(t, err, errors.New("invalid DSN: missing '@' or '(' or ')' separating the necessary parts"))
		})
		Convey("OpenConnector fail,change fail,invalid time unit suffix\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout=10xxs", info.Username, info.Password, info.Host, info.Port)
			_, err := OpenConnector(dsn)
			assert.Equal(t, err, errors.New("time: unknown unit \"xxs\" in duration \"10xxs\""))
		})
		Convey("OpenConnector success,test invalid param continue\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout", info.Username, info.Password, info.Host, info.Port)
			_, err := OpenConnector(dsn)
			assert.Equal(t, err, nil)
		})
		Convey("OpenConnector success,case param two\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?readTimeout=10s&timeout=10s", info.Username, info.Password, info.Host, info.Port)
			_, err := OpenConnector(dsn)
			assert.Equal(t, err, nil)
		})
		Convey("OpenConnector success\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout=10s&autocommit=true&readTimeout=10s", info.Username, info.Password, info.Host, info.Port)
			_, err := OpenConnector(dsn)
			assert.Equal(t, err, nil)
		})
	})
}

func TestExecContext(t *testing.T) {
	Convey("Test dmdb.ExecContext\n", t, func() {
		info := getTestDBInfo()
		Convey("ExecContext fail\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout=10s&autocommit=true&readTimeout=10s", info.Username, info.Password, info.Host, info.Port)
			dmConn, err := Open(dsn)
			assert.Equal(t, err, nil)
			_, err = dmConn.(driver.ExecerContext).ExecContext(context.Background(), "insert t1(id) values 1", nil)
			assert.NotEqual(t, err, nil)
		})
		Convey("ExecContext success\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout=10s&autocommit=true&readTimeout=10s", info.Username, info.Password, info.Host, info.Port)
			dmConn, err := Open(dsn)
			assert.Equal(t, err, nil)
			_, err = dmConn.(driver.ExecerContext).ExecContext(context.Background(), "CREATE TABLE IF NOT EXISTS t1(id int)", []driver.NamedValue{})
			assert.Equal(t, err, nil)
			_, err = dmConn.(driver.ExecerContext).ExecContext(context.Background(), "insert t1(id) values (1)", []driver.NamedValue{})
			assert.Equal(t, err, nil)
			_, err = dmConn.(driver.ExecerContext).ExecContext(context.Background(), "DROP TABLE IF EXISTS t1", []driver.NamedValue{})
			assert.Equal(t, err, nil)
		})
	})
}

func TestQueryContext(t *testing.T) {
	Convey("Test dmdb.QueryContext\n", t, func() {
		info := getTestDBInfo()
		Convey("QueryContext fail\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout=10s&autocommit=true&readTimeout=10s", info.Username, info.Password, info.Host, info.Port)
			dmConn, err := Open(dsn)
			assert.Equal(t, err, nil)
			_, err = dmConn.(driver.QueryerContext).QueryContext(context.Background(), "selectt 1", nil)
			assert.NotEqual(t, err, nil)
		})
		Convey("QueryContext success\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout=10s&autocommit=true&readTimeout=10s", info.Username, info.Password, info.Host, info.Port)
			dmConn, err := Open(dsn)
			assert.Equal(t, err, nil)
			_, err = dmConn.(driver.QueryerContext).QueryContext(context.Background(), "select 1", nil)
			assert.Equal(t, err, nil)
		})
	})
}

func TestPrepareContext(t *testing.T) {
	Convey("Test dmdb.PrepareContext\n", t, func() {
		info := getTestDBInfo()
		Convey("PrepareContext fail,no slash\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout=10s&autocommit=true&readTimeout=10s", info.Username, info.Password, info.Host, info.Port)
			dmConn, err := Open(dsn)
			assert.Equal(t, err, nil)
			_, err = dmConn.(driver.ConnPrepareContext).PrepareContext(context.Background(), "selectt 1")
			assert.NotEqual(t, err, nil)
		})
		Convey("PrepareContext success\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout=10s&autocommit=true&readTimeout=10s", info.Username, info.Password, info.Host, info.Port)
			dmConn, err := Open(dsn)
			assert.Equal(t, err, nil)
			_, err = dmConn.(driver.ConnPrepareContext).PrepareContext(context.Background(), "select 1")
			assert.Equal(t, err, nil)
		})
	})
}
