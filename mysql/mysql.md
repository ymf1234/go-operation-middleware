##连接

---

> ### 下载依赖
```shell
go get -u github.com/go-sql-driver/mysql
```

> ###使用MySQL驱动
```go
func Open(driverName, dataSourceName string) (*DB, error)
```

Open打开一个dirverName指定的数据库，dataSourceName指定数据源，一般至少包括数据库文件名和其它连接必要的信息。

```go
package main
import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
   // DSN:Data Source Name
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()  // 注意这行代码要写在上面err判断的下面
}
```

> 初始化连接

Open函数可能只是验证其参数格式是否正确，实际上并不创建与数据库的连接。如果要检查数据源的名称是否真实有效，应该调用Ping方法。

返回的DB对象可以安全地被多个goroutine并发使用，并且维护其自己的空闲连接池。因此，Open函数应该仅被调用一次，很少需要关闭这个DB对象。

接下来，我们定义一个全局变量db，用来保存数据库连接对象。将上面的示例代码拆分出一个独立的initDB函数，只需要在程序启动时调用一次该函数完成全局变量db的初始化，其他函数中就可以直接使用全局变量db了。`（注意下方的注意）`
```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// 定义存储db全局变量
var db *sql.DB

// 初始化数据库
func initDB() (err error) {
	// 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/book?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err = sql.Open("mysql", dsn)

	if err != nil {
		return err
	}
	defer db.Close() // 断开连接
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()

	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := initDB() // 调用输出化数据库的函数
	if err != nil {
		fmt.Printf("初始化失败,err:%v\n", err)
		return
	}

	fmt.Printf("连接成功\n")
	return
}
```
其中sql.DB是表示连接的数据库对象（结构体实例），它保存了连接数据库相关的所有信息。它内部维护着一个具有零到多个底层连接的连接池，它可以安全地被多个goroutine同时使用。

> ###SetMaxOpenConns
```go
func (db *DB) SetMaxOpenConns(n int)
```
`SetMaxOpenConns`设置与数据库建立连接的最大数目。 如果n大于0且小于最大闲置连接数，会将最大闲置连接数减小到匹配最大开启连接数的限制。 如果n<=0，不会限制最大开启连接数，默认为0（无限制）。
> ###SetMaxIdleConns
```go
func (db *DB) SetMaxIdleConns(n int)
```
`SetMaxIdleConns`设置连接池中的最大闲置连接数。 如果n大于最大开启连接数，则新的最大闲置连接数会减小到匹配最大开启连接数的限制。 如果n<=0，不会保留闲置连接。

##CRUD

---

>建库建表
我们先在MySQL中创建一个名为test的数据库
```mysql
CREATE DATABASE test;

use test;

-- 创建测试表
CREATE TABLE `user` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(10) NOT NULL COMMENT '用户名',
    `password` varchar(15) NOT NULL DEFAULT '123456' COMMENT '密码',
    `address` varchar(25) DEFAULT NULL COMMENT '地址',
    `phone` varchar(15) DEFAULT NULL COMMENT '手机号',
    `money` int(11) NOT NULL DEFAULT '0' COMMENT '钱 单位:分',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
```

> 查询
```sql
    
```