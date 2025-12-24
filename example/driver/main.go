package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AISHU-Technology/proton-rds-sdk-go/driver"
	"github.com/AISHU-Technology/proton-rds-sdk-go/example/driver/kdb"
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
func testDM8() {
	os.Setenv("DB_TYPE", "dm8")
	info := getTestDBInfo()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/deploy?timeout=10s&client_encoding=utf-8", info.Username, info.Password, info.Host, info.Port)
	op, err := sql.Open("proton-rds", dsn)

	// os.Setenv("DB_TYPE", "kdb9")
	// op, err := sql.Open("proton-rds", "system:system@tcp(localhost:54321)/anyshare?timeout=1s")

	//op, err := sql.Open("proton-rds", "root:eisoo.com123@tcp(localhost:3306)/test?timeout=10s")
	if err != nil {
		fmt.Println(err)
		return
	}
	op.Exec("drop table  t111")
	_, err = op.Exec("create table if not exists t111(`id` int,`time` datetime(0), `name` varchar)")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = op.Exec("insert into t111 values(?,?,?)", 1, time.Now(), "克路")
	if err != nil {
		fmt.Println(err)
		return
	}

	var t driver.Time
	var n string
	//var t time.Time

	rows, err := op.Query("select `time`,`name` from t111 where  `id`=?", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&t, &n)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(t, len(t.String()), n)
	}
}

func main() {
	testDM8()
	kdb.Test()
}
