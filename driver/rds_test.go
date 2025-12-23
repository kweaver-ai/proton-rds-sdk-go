package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

var dbType []string = []string{"MARIADB", "DM8", "KDB9"}

type TestDBInfo struct {
	Host     string
	Port     int
	Username string
	Password string
}

func getTestDBInfo(db_type string) TestDBInfo {
	db_type = strings.ToUpper(db_type)
	switch db_type {
	case "MARIADB":
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
	case "DM8":
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
	case "KDB9":
		user := os.Getenv("KDB_TEST_USER")                    // 例如设置为 "test_user"
		password := os.Getenv("KDB_TEST_PASSWORD")            // 例如设置为 "test_pwd_123"
		host := os.Getenv("KDB_TEST_HOST")                    // 例如设置为 "localhost"
		port, err := strconv.Atoi(os.Getenv("KDB_TEST_PORT")) // 例如设置为 "3306"
		if err != nil {
			log.Fatalf("KDB_TEST_PORT is not a number: %v", err)
		}
		return TestDBInfo{
			Host:     host,
			Port:     port,
			Username: user,
			Password: password,
		}
	default:
		log.Fatalf("DB_TYPE %s is not supported", db_type)
	}
	return TestDBInfo{}
}

func TestOpen1(t *testing.T) {
	err := os.Setenv("DB_TYPE", dbType[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Convey("Test db.Open\n", t, func() {
		info := getTestDBInfo(dbType[0])
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?timeout=10s", info.Username, info.Password, info.Host, info.Port)
		Convey("Error dsn, return error\n", func() {
			_, err = RDSDriver{}.Open("xxxxx")
			assert.NotNil(t, err)
		})
		Convey("Open fail\n", func() {
			_, err := RDSDriver{}.Open("/test?timeout=10s&readTimeout=10s")
			assert.NotNil(t, err)
		})
		Convey("Open fail,invalid time param\n", func() {
			errdsn := dsn + "&timeout=10xxxxxs"
			_, err := RDSDriver{}.Open(errdsn)
			assert.NotNil(t, err)
		})
		Convey("Open ok\n", func() {
			_, err := RDSDriver{}.Open(dsn)
			assert.Nil(t, err)
		})
	})
}

func TestOpen2(t *testing.T) {
	err := os.Setenv("DB_TYPE", dbType[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Convey("Test db.Open\n", t, func() {
		info := getTestDBInfo(dbType[1])
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout=10s&readTimeout=10s", info.Username, info.Password, info.Host, info.Port)
		Convey("Error dsn, return error\n", func() {
			_, err = RDSDriver{}.Open("xxxxx")
			assert.NotNil(t, err)
		})
		Convey("Open ok\n", func() {
			_, err := RDSDriver{}.Open(dsn)
			assert.Nil(t, err)
		})
	})
}

func TestOpen3(t *testing.T) {
	err := os.Setenv("DB_TYPE", dbType[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Convey("Test db.Open\n", t, func() {
		info := getTestDBInfo(dbType[2])
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/proton?timeout=10s", info.Username, info.Password, info.Host, info.Port)
		Convey("Error dsn, return error\n", func() {
			_, err = RDSDriver{}.Open("xxxxx")
			assert.NotNil(t, err)
		})
		Convey("Open ok\n", func() {
			_, err := RDSDriver{}.Open(dsn)
			assert.Nil(t, err)
		})
	})
}

func TestInit(t *testing.T) {
	err := os.Setenv("DB_TYPE", "DM8")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Convey("Test init\n", t, func() {
		info := getTestDBInfo(dbType[1])
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout=10s&readTimeout=10s", info.Username, info.Password, info.Host, info.Port)
		Convey("Error driver name, return error\n", func() {
			_, err := sql.Open("xxx", dsn)
			assert.NotNil(t, err)
		})
		Convey("proton-rds driver\n", func() {
			_, err := sql.Open("proton-rds", dsn)
			assert.Nil(t, err)
		})
	})
}

func TestOpenConnector1(t *testing.T) {
	err := os.Setenv("DB_TYPE", dbType[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Convey("Test db.OpenConnector\n", t, func() {
		info := getTestDBInfo(dbType[0])
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?timeout=10s", info.Username, info.Password, info.Host, info.Port)
		Convey("Error dsn, return error\n", func() {
			_, err = RDSDriver{}.OpenConnector("xxxxx")
			assert.NotNil(t, err)
		})
		Convey("OpenConnector fail\n", func() {
			_, err := RDSDriver{}.OpenConnector("test?timeout=10s&readTimeout=10s")
			assert.NotNil(t, err)
		})
		Convey("OpenConnector fail,invalid time param\n", func() {
			errdsn := dsn + "&timeout=10xxxxxs"
			_, err := RDSDriver{}.OpenConnector(errdsn)
			assert.NotNil(t, err)
		})
		Convey("OpenConnector ok\n", func() {
			_, err := RDSDriver{}.OpenConnector(dsn)
			assert.Nil(t, err)
		})
	})
}

func TestOpenConnector2(t *testing.T) {
	err := os.Setenv("DB_TYPE", dbType[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Convey("Test db.OpenConnector\n", t, func() {
		info := getTestDBInfo(dbType[1])
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/SYSDBA?timeout=10s&readTimeout=10s", info.Username, info.Password, info.Host, info.Port)
		Convey("Error dsn, return error\n", func() {
			_, err = RDSDriver{}.OpenConnector("xxxxx")
			assert.NotNil(t, err)
		})
		Convey("OpenConnector ok\n", func() {
			_, err := RDSDriver{}.OpenConnector(dsn)
			assert.Nil(t, err)
		})
	})
}

func TestOpenConnector3(t *testing.T) {
	err := os.Setenv("DB_TYPE", dbType[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Convey("Test db.OpenConnector\n", t, func() {
		info := getTestDBInfo(dbType[2])
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/proton?timeout=10s", info.Username, info.Password, info.Host, info.Port)
		Convey("Error dsn, return error\n", func() {
			_, err = RDSDriver{}.OpenConnector("xxxxx")
			assert.NotNil(t, err)
		})
		Convey("OpenConnector ok\n", func() {
			_, err := RDSDriver{}.OpenConnector(dsn)
			assert.Nil(t, err)
		})
	})
}
