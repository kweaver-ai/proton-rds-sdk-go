package kdb

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func TestTime(op *sql.DB) {
	tableSql := "CREATE TABLE IF NOT EXISTS `test_time`(" +
		"`id` INT," +
		"`time` DATETIME" +
		")"
	_, err := op.Exec(tableSql)
	if err != nil {
		fmt.Println(err)
		return
	}

	type st struct {
		id int64
		t  time.Time
	}

	os := st{
		id: 1,
		t:  time.Now(),
	}
	_, err = op.Exec("INSERT INTO `test_time` VALUES(?,?)", os.id, os.t)
	if err != nil {
		fmt.Println(err)
		return
	}

	row := op.QueryRow("SELECT `id`, `time` FROM `test_time` WHERE id=?", 1)
	ns := st{}
	err = row.Scan(
		&ns.id,
		&ns.t,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	newTimeStr := ns.t.Format(time.RFC3339)
	orgTimeStr := os.t.Format(time.RFC3339)
	if ns.id != os.id ||
		newTimeStr != orgTimeStr {
		log.Fatalf("data not match: new: %v, org: %v", ns, os)
	}

	fmt.Println("success")
}
