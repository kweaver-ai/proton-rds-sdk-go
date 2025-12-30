package main

import (
	"fmt"

	_ "github.com/kweaver-ai/proton-rds-sdk-go/driver"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
)

func main() {
	// drivername是数据库驱动注册的名字(https://pkg.go.dev/database/sql#Register)
	// dataSourceName是master节点和backup节点的dsn逗号分割拼接
	// create op once
	connInfo := sqlx.DBConfig{}
	connInfo.User = "username"
	connInfo.Password = "password"
	connInfo.Host = "localhost"
	connInfo.Port = 5236
	connInfo.HostRead = "localhost"
	connInfo.PortRead = 5236
	connInfo.Database = "deploy"
	connInfo.ParseTime = "true"
	connInfo.Loc = "Local"
	op, err := sqlx.NewDB(&connInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	// op.Exec: exec sql on master node. Use it like DB.Exec
	_, err = op.Exec("create table if not exists t1(id int, name varchar(20) default null)")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = op.Exec("truncate table t1")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = op.Exec("insert into t1(id) values(?)", 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = op.Exec("insert into t1(id) values(?)", 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = op.Exec("insert into t1(id) values(?)", 3)
	if err != nil {
		fmt.Println(err)
		return
	}

	var v int

	// exec sql on backup node. Use it like DB.QueryRow
	row := op.QueryRow("select max(id) from t1")
	err = row.Scan(&v)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("max id: ", v)

	// exec sql on backup node. Use it like DB.Query
	rows, err := op.Query("select id from t1 where id > ?", 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&v)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("id: ", v)
	}
}
