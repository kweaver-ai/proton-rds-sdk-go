package kdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
)

type TestDBInfo struct {
	Host     string
	Port     int
	Username string
	Password string
}

func getTestDBInfo() TestDBInfo {
	user := os.Getenv("KDB9_TEST_USER")                    // 例如设置为 "test_user"
	password := os.Getenv("KDB9_TEST_PASSWORD")            // 例如设置为 "test_pwd_123"
	host := os.Getenv("KDB9_TEST_HOST")                    // 例如设置为 "localhost"
	port, err := strconv.Atoi(os.Getenv("KDB9_TEST_PORT")) // 例如设置为 "3306"
	if err != nil {
		log.Fatalf("KDB9_TEST_PORT is not a number: %v", err)
	}
	return TestDBInfo{
		Host:     host,
		Port:     port,
		Username: user,
		Password: password,
	}
}

func Test() {
	os.Setenv("DB_TYPE", "kdb9")
	info := getTestDBInfo()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?timeout=10s", info.Username, info.Password, info.Host, info.Port)
	op, err := sql.Open("proton-rds", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = op.Exec("DROP SCHEMA IF EXISTS `test` CASCADE")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = op.Exec("CREATE SCHEMA IF NOT EXISTS `test`")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = op.Exec("SET SEARCH_PATH TO `test`")
	if err != nil {
		fmt.Println(err)
		return
	}

	TestTime(op)
	TestBlob(op)
	TestHydraSelect(op)
}
