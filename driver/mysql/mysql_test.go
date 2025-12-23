package mysql

import (
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
	user := os.Getenv("MYSQL_TEST_USER")                    // 例如设置为 "test_user"
	password := os.Getenv("MYSQL_TEST_PASSWORD")            // 例如设置为 "test_pwd_123"
	host := os.Getenv("MYSQL_TEST_HOST")                    // 例如设置为 "localhost"
	port, err := strconv.Atoi(os.Getenv("MYSQL_TEST_PORT")) // 例如设置为 "3306"
	if err != nil {
		log.Fatalf("MYSQL_TEST_PORT is not a number: %v", err)
	}
	return TestDBInfo{
		Host:     host,
		Port:     port,
		Username: user,
		Password: password,
	}
}

func TestOpen(t *testing.T) {
	Convey("Test mysql.Open\n", t, func() {
		info := getTestDBInfo()
		Convey("Open fail\n", func() {
			error_port := 30036
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", info.Username, info.Password, info.Host, error_port)
			_, err := Open(dsn)
			assert.NotNil(t, err)
		})
		Convey("Open success\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", info.Username, info.Password, info.Host, info.Port)
			_, err := Open(dsn)
			assert.Nil(t, err)
		})
	})
}

func TestOpenConnector(t *testing.T) {
	Convey("Test mysql.OpenConnector\n", t, func() {
		info := getTestDBInfo()
		Convey("Open fail\n", func() {
			error_port := 30036
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)test", info.Username, info.Password, info.Host, error_port)
			_, err := OpenConnector(dsn)
			assert.NotNil(t, err)
		})
		Convey("OpenConnector success\n", func() {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", info.Username, info.Password, info.Host, info.Port)
			_, err := OpenConnector(dsn)
			assert.Nil(t, err)
		})
	})
}
