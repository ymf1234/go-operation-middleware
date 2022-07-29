package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// 定义存储db全局变量
var db *sql.DB

// 初始化数据库
// setMaxOpenConns  设置与数据库建立连接的最大数目。 如果n大于0且小于最大闲置连接数，会将最大闲置连接数减小到匹配最大开启连接数的限制。 如果n<=0，不会限制最大开启连接数，默认为0（无限制）。
// setMaxIdleConns  设置连接池中的最大闲置连接数。 如果n大于最大开启连接数，则新的最大闲置连接数会减小到匹配最大开启连接数的限制。 如果n<=0，不会保留闲置连接。
func initDB(setMaxOpenConns int, setMaxIdleConns int) (err error) {
	// 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err = sql.Open("mysql", dsn)

	if err != nil {
		return err
	}
    //defer db.Close() // 断开连接
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()

	if err != nil {
		return err
	}
	// 设置与数据库建立连接的最大数目
	db.SetMaxOpenConns(setMaxOpenConns)
	//设置连接池中的最大闲置连接数
	db.SetMaxIdleConns(setMaxIdleConns)
	return nil
}

func main() {
	err := initDB(50, 10) // 调用输出化数据库的函数
	if err != nil {
		fmt.Printf("初始化失败,err:%v\n", err)
		return
	}

	fmt.Printf("连接成功\n")

	queryRowDemo()

	return
}


// CRUD

type user struct {
	id int
	name string
	password string
	address string
	phone string
	money int
}


// 查询单条数据示例
func queryRowDemo() {
	sqlStr := "select id, name, password, address, phone,money  from `user` where id = ?"
	var u user
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	row := db.QueryRow(sqlStr, 1)
	err := row.Scan(&u.id, &u.name, &u.password, &u.address, &u.phone, &u.money)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Printf("date: %+v\n", u)
}

// 多行查询
func queryMultiRowDemo()  {
	sqlStr := "select id, name, password, address, phone,money  from `user` where id > ?"
	var u user
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		
	}



}
