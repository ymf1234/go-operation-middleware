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

	/**
	CRUD
	 */
	// 查询一条
	//fmt.Println("查询单条")
	//queryRowDemo()

	// 查询多条
	//fmt.Println("查询多条")
	//queryMultiRowDemo()

	// 插入数据
	//fmt.Println("插入")
	//insertRowDemo()

	// 修改数据
	//fmt.Println("修改")
	//updateRowDemo()

	// 删除数据
	//deleteRowDemo()



	/**
	预处理
	*/
	//fmt.Println("预处理查询")
	//prepareQueryDemo()

	//fmt.Println("预处理新增")
	//prepareInsertDemo()

	// SQL注入
	//fmt.Println("SQL注入")
	//sqlInjectDemo("xxx' or 1=1#")
	//fmt.Println()
	//sqlInjectDemo("xxx' union select * from user #")
	//fmt.Println()
	//sqlInjectDemo("xxx' or (select count(*) from user) <10 #")


	// 事务
	transactionDemo()
	db.Close()
	return
}

// CRUD

type user struct {
	id       int
	name     string
	password string
	address  string
	phone    string
	money    int
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
func queryMultiRowDemo() {
	sqlStr := "select id, name, password, address, phone,money  from `user` where id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.password, &u.address, &u.phone, &u.money)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("data: %+v\n", u)
	}
}

// 插入数据
func insertRowDemo() {
	sqlStr := "insert into user(name,password,address,phone,money) value(?, ?, ?, ?, ?)"
	exec, err := db.Exec(sqlStr, "小福", "123456", "宇宙", "17676767676", "100")
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	id, err := exec.LastInsertId() // 新插入数据的id

	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", id)
}

// 更新数据
func updateRowDemo() {
	sqlStr := "update user set money = ? where id = ?"
	exec, err := db.Exec(sqlStr, 20, 3)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	n,err := exec.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

// 删除数据
func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	exec, err := db.Exec(sqlStr, 6)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err := exec.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}


// 预处理

// 预处理查询
func prepareQueryDemo() {
	sqlStr := "select id, name, password, address, phone,money  from `user` where id > ?"
	prepare, err := db.Prepare(sqlStr)

	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}
	defer prepare.Close()

	rows, err := prepare.Query(0)

	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}

	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.password, &u.address, &u.phone, &u.money)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}

		fmt.Printf("data: %+v\n", u)
	}
}

// 预处理插入
func prepareInsertDemo() {
	sqlStr := "insert into user(name,password,address,phone,money) value(?, ?, ?, ?, ?)"
	prepare, err := db.Prepare(sqlStr)

	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}

	defer prepare.Close()

	result1, err := prepare.Exec("六便士", "123456", "书", "12345678900", "34")
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	id1, err := result1.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", id1)

	result2, err := prepare.Exec("六便士1", "123456", "书", "12345678900", "34")
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	id2, err := result2.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", id2)

}


// SQL注入
func sqlInjectDemo(name string) {
	sqlStr := fmt.Sprintf("select id, name, password, address, phone,money from user where name='%s'", name)
	fmt.Printf("SQL:%s\n", sqlStr)
	var u user
	err := db.QueryRow(sqlStr).Scan(&u.id, &u.name, &u.password, &u.address, &u.phone, &u.money)

	if err != nil {
		fmt.Printf("exec failed, err:%v\n", err)
		return
	}
	fmt.Printf("user:%#v\n", u)
}

// 事务
func transactionDemo() {
	tx, err := db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		fmt.Printf("事务开启失败, err:%v\n", err)
		return
	}

	sqlStr := "update user set money = 50 where id = ?"
	ret1, err := tx.Exec(sqlStr, 2)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("执行sql1失败, err:%v\n", err)
	}

	affRow1, err := ret1.RowsAffected() // 操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("执行 ret1.RowsAffected() 失败, err:%v\n", err)
		return
	}
	sqlStr2 := "update user set money = 50 where id = ?"
	ret2, err := tx.Exec(sqlStr2, 3)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("执行sql1失败, err:%v\n", err)
	}

	affRow2, err := ret2.RowsAffected() // 操作影响的行数
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("执行 ret2.RowsAffected() 失败, err:%v\n", err)
		return
	}

	fmt.Println(affRow1, affRow2)
	if affRow1 == 1 && affRow2 == 1 {
		fmt.Println("事务提交啦...")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		fmt.Println("事务回滚啦...")
	}

	fmt.Println("执行事务成功")
}