package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var db *sqlx.DB

func initDB() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return nil
}

func main() {
	err := initDB()
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	defer db.Close()
	//fmt.Println("单条查询")
	//queryRowDemo()
	//fmt.Println("多条查询")
	//queryMultiRowDemo()
	//
	//fmt.Println("插入")
	//insertRowDemo()
	//
	//fmt.Println("修改")
	//updateRowDemo()
	//
	//fmt.Println("删除")
	//deleteRowDemo()
	//
	//fmt.Println("预处理")
	//nameQuery()

	//fmt.Println("事务")
	//err = transactionDemo()
	//fmt.Println(err)

	u1 := Users{Name: "七米", Password: "123456", Address: "123456", Phone: "1333333333", Money: 20}
	u2 := Users{Name: "q1mi", Password: "123456", Address: "123456", Phone: "1333333333", Money: 20}
	u3 := Users{Name: "小王子", Password: "123456", Address: "123456", Phone: "1333333333", Money: 20}

	// 方法1
	//users := []*Users{&u1, &u2, &u3}
	//err = BatchInsertUsers(users)
	//if err != nil {
	//	fmt.Printf("BatchInsertUsers failed, err:%v\n", err)
	//}

	// 方法2
	//fmt.Println("BatchInsertUsers2")
	//users2 := []interface{}{u1, u2, u3}
	//err = BatchInsertUsers2(users2)
	//if err != nil {
	//	fmt.Printf("BatchInsertUsers2 failed, err:%v\n", err)
	//}

	// 方法3
	users3 := []*Users{&u1, &u2, &u3}
	err = BatchInsertUsers3(users3)
	if err != nil {
		fmt.Printf("BatchInsertUsers3 failed, err:%v\n", err)
	}
}


type user struct {
	Id       int
	Name     string
	Password string
	Address  string
	Phone    string
	Money    int
}

// 单条查询
func queryRowDemo() {
	sqlStr := "select id, name, password, address, phone,money  from `user` where id = ?"
	var u user
	err := db.Get(&u, sqlStr, 1)

	if err != nil {
		fmt.Printf("获取数据失败，err:%v\n", err)
		return
	}

	fmt.Printf("date: %+v\n", u)
}

// 多条查询
func queryMultiRowDemo() {
	sqlStr := "select id, name, password, address, phone,money  from `user` where id > ?"
	var users []user
	err := db.Select(&users, sqlStr, 1)

	if err != nil {
		fmt.Printf("获取数据失败，err:%v\n", err)
		return
	}

	fmt.Printf("users: %+v\n", users)
}

// 插入
func insertRowDemo() {
	sqlStr := "insert into user(name,password,address,phone,money) value(?, ?, ?, ?, ?)"
	ret, err := db.Exec(sqlStr, "小福1", "123456", "宇宙", "17676767676", "100")
	if err != nil {
		fmt.Printf("新增失败, err:%v\n", err)
		return
	}

	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("获取id失败, err:%v\n", err)
		return
	}
	fmt.Printf("插入成功, id: %d.\n", theID)
}

// 更新数据
func updateRowDemo() {
	sqlStr := "update user set money = ? where id = ?"
	ret, err := db.Exec(sqlStr, 300, 9)
	if err != nil {
		fmt.Printf("修改失败， err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("获取 RowsAffected 失败, err:%v\n", err)
		return
	}
	fmt.Printf("修改成功, 影响行数: %d\n", n)
}

// 删除数据
func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := db.Exec(sqlStr, 8)
	if err != nil {
		fmt.Printf("删除失败， err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("获取 RowsAffected 失败, err:%v\n", err)
		return
	}
	fmt.Printf("删除成功, 影响行数: %d\n", n)
}

// 预处理
func nameQuery() {
	sqlStr := "select * from user where name=:name"
	// 使用map做命名查询
	rows, err := db.NamedQuery(sqlStr, map[string]interface{}{"name": "小福1"})
	if err != nil {
		fmt.Printf("db.NamedQuery failed, err:%v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}

	u := user{
		Name: "六便士",
	}
	// 使用结构体命名查询，根据结构体字段的 db tag进行映射
	rows, err = db.NamedQuery(sqlStr, u)
	if err != nil {
		fmt.Printf("db.NamedQuery failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next(){
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}
}

// 事务
func transactionDemo() (err error){
	tx, err := db.Beginx() // 开启事务
	if err != nil {
		fmt.Printf("开启事务失败, err:%v\n", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(err) // panic 之后 回滚
		} else if err != nil {
			fmt.Println("rollback")
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
			fmt.Println("commit")
		}
	}()


	sqlStr1 := "Update user set money=20 where id=?"
	rs, err := tx.Exec(sqlStr1, 1)
	if err!= nil{
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("执行 sqlStr1 失败")
	}

	sqlStr2 := "Update user set money=50 where id=?"
	rs, err = tx.Exec(sqlStr2, 5)
	if err!=nil{
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	return err
}

// sqlx.In的批量插入示例

type Users struct {
	Name     string `db:"name"`
	Password string `db:"password"`
	Address  string `db:"address"`
	Phone    string `db:"phone"`
	Money    int `db:"money"`
}

// 自己拼接语句实现批量插入
func BatchInsertUsers(users []*Users) error {
	// 存放 (?, ?) 的slice
	valueStrings := make([]string, 0, len(users))

	// 存放values的slice
	valueArgs := make([]interface{}, 0, len(users) * 2)

	// 遍历users准备相关数据
	for _, u := range users {
		// 此处占位符要与插入值的个数对应
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, u.Name)
		valueArgs = append(valueArgs, u.Password)
		valueArgs = append(valueArgs, u.Address)
		valueArgs = append(valueArgs, u.Phone)
		valueArgs = append(valueArgs, u.Money)
	}

	// 自行拼接要执行的具体语句
	stmt := fmt.Sprintf("INSERT INTO user (name,password,address,phone,money) VALUES %s",
		strings.Join(valueStrings, ","))
	_, err := db.Exec(stmt, valueArgs...)
	return err
}


// 使用sqlx.In实现批量插入
func (u Users) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Password, u.Address, u.Phone, u.Money}, nil
}

// BatchInsertUsers2 使用sqlx.In帮我们拼接语句和参数, 注意传入的参数是[]interface{}
func BatchInsertUsers2(users []interface{}) error {
	fmt.Println("users: ", users)
	query, args, _ := sqlx.In(
		"INSERT INTO user (name,password,address,phone,money) VALUES (?), (?), (?)",
		users..., // 如果arg实现了 driver.Valuer, sqlx.In 会通过调用 Value()来展开它
	)
	fmt.Println(query) // 查看生成的querystring
	fmt.Println(args)  // 查看生成的args
	_, err := db.Exec(query, args...)
	return err
}

// BatchInsertUsers3 使用NamedExec实现批量插入
func BatchInsertUsers3(users []*Users) error {
	_, err := db.NamedExec("INSERT INTO user (name,password,address,phone,money) VALUES (:name,:password,:address,:phone,:money)", users)
	return err
}

// sqlx.In的查询示例
// QueryByIDs 根据给定ID查询
func QueryByIDs(ids []int)(users []Users, err error){
	// 动态填充id
	query, args, err := sqlx.In("SELECT name FROM user WHERE id IN (?)", ids)
	if err != nil {
		return
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)

	err = db.Select(&users, query, args...)
	return
}

// QueryAndOrderByIDs 按照指定id查询并维护顺序
func QueryAndOrderByIDs(ids []int)(users []Users, err error){
	// 动态填充id
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, fmt.Sprintf("%d", id))
	}
	query, args, err := sqlx.In("SELECT name, age FROM user WHERE id IN (?) ORDER BY FIND_IN_SET(id, ?)", ids, strings.Join(strIDs, ","))
	if err != nil {
		return
	}

	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)

	err = db.Select(&users, query, args...)
	return
}