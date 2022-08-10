package main

// https://learnku.com/docs/gorm/v2
// https://gorm.io/zh_CN/docs/

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type sqliteModel struct {
	db *gorm.DB
}

// 初始化sqlite
/*func (s *sqliteModel) sqliteInit() (err error) {
	s.db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		fmt.Printf("sqlite 连接失败, err: %v\n", err)
	}
	return err
}*/

// 迁移 schema
func (s *sqliteModel) sqliteAutoMigrate(i ...interface{}) {
	s.db.AutoMigrate(i...)
}

func (s *sqliteModel) sqliteCreate(i interface{}) {
	s.db.Create(i)
}



/**
	MySQL
 */
type mysqlModel struct {
	db *gorm.DB
	err error
}

type GormUser struct {
	ID           uint
	Name         string
	//Email        *string
	Age          uint8
	Birthday     time.Time
	//MemberNumber sql.NullString
	//ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (GormUser) TableName() string {
	return "gorm_users"
}

func (m *mysqlModel) mysqlInit() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	dialector := mysql.New(mysql.Config{
		DSN:                       dsn, // DSN data source name
		DefaultStringSize:         256, // string 类型字段的默认长度
		DisableDatetimePrecision:  true,// 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,// 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,// 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,// 根据当前 MySQL 版本自动配置
	})
	//m.db, m.err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	m.db, m.err = gorm.Open(dialector, &gorm.Config{})

	if m.err != nil {
		fmt.Printf("sqlite 连接失败, err: %v\n", m.err)
		panic(m.err)
	}
}

func (m *mysqlModel) mysqlCreate() {
	user := GormUser{
		Name:     "Go ORM",
		Age:      7,
		Birthday: time.Now(),
	}

	//result := m.db.Create(&user) // 通过数据的指针来创建

	//result := m.db.Select("name", "age", "created_at").Create(&user) // 用指定的字段创建记录
	result := m.db.Omit("name", "age", "created_at").Create(&user) // 创建一个记录且一同忽略传递给略去的字段值。
	fmt.Println("插入数据的主键:", user.ID)
	fmt.Println("返回 error:", result.Error)
	fmt.Println("返回插入记录的条数:", result.RowsAffected )
}

// 批量插入
func (m *mysqlModel) mysqlBCreate() {
	var users = []GormUser{
		{Name: "Go ORM1"},
		{Name: "Go ORM2"},
		{Name: "Go ORM3"},
	}

	result := m.db.Create(&users) // 通过数据的指针来创建

	// 数量为 2
	//result := m.db.CreateInBatches(&users, 2) // 分批插入

	for _, user := range users {
		fmt.Println("插入数据的主键:", user.ID)
	}



	fmt.Println("返回 error:", result.Error)
	fmt.Println("返回插入记录的条数:", result.RowsAffected )
}

// 查询
func (m *mysqlModel) mysqlSelect() {
	var u GormUser
	// 获取第一条记录（主键升序）
	//result := m.db.First(&u) // SELECT * FROM gorm_users ORDER BY id LIMIT 1;
	// 获取一条记录，没有指定排序字段
	//result := m.db.Take(&u)

	// 获取最后一条记录（主键降序）
	result := m.db.Last(&u) // SELECT * FROM users ORDER BY id DESC LIMIT 1

	fmt.Printf("查询数据：%v\n", u)
	fmt.Println(result.RowsAffected) // 返回找到的记录数
	// 检查 ErrRecordNotFound 错误
	is := errors.Is(result.Error, gorm.ErrRecordNotFound)
	fmt.Println(result.Error, is)        // returns error or nil

}

func (m *mysqlModel) mysqlUpdate() {
	var u GormUser
	//m.db.Model(&u).Where("name = ?", "Go ORM").Update("name", "hello1")

	m.db.Model(&u).Where("name = ?", "Go ORM").Updates(GormUser{Name: "A", Age: 20})
}


func (m *mysqlModel) mysqlDelete() {
	var u GormUser
	//m.db.Model(&u).Where("name = ?", "Go ORM").Update("name", "hello1")

	m.db.Delete(&u, 10)
}


func main() {
	/*sqlite1 := &sqliteModel{}
	_ = sqlite1.sqliteInit()
	sqlite1.sqliteAutoMigrate(&Product{})
	sqlite1.sqliteCreate(&Product{Code: "D42", Price: 100})*/

	gorm1 := &mysqlModel{}
	gorm1.mysqlInit()
	//gorm1.mysqlCreate()
	//gorm1.mysqlBCreate()
	//gorm1.mysqlSelect()
	//gorm1.mysqlUpdate()
	gorm1.mysqlDelete()
}
