package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var db *sql.DB

func InitMysql() {
	log.Info("InitMysql....")
	db, _ = sql.Open("mysql", "root:2220718-lzqj@tcp(127.0.0.1:3306)/homecloud")
	//if err != nil { //判断成功失败
	//	fmt.Println("数据库连接失败")
	//	panic(err.Error()) // proper error handling instead of panic in your app
	//	//如果连接失败 先检查数据库服务是否已经启动 net start mysql启动数据库
	//	return
	//}
	//err = db.Ping()
	//if err != nil {
	//	panic(err.Error()) // proper error handling instead of panic in your app
	//}
	CreateTableWithUser()
	CreateTableWithFile()
}

//创建用户表
func CreateTableWithUser() {
	log.Info("CreateTableWithUser")
	sql := `CREATE TABLE IF NOT EXISTS users(
        id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,
        username VARCHAR(64),
        password VARCHAR(64),
        status INT(4),
        createtime INT(10)
        );`

	ModifyDB(sql)
}

func ModifyDB(sql string, args ...interface{}) (int64, error) {
	log.Info("modifyDB")
	result, err := db.Exec(sql, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return count, nil
}

//查询
func QueryRowDB(sql string) *sql.Row{
	return db.QueryRow(sql)
}

func QueryDB(sql string) (*sql.Rows, error) {
	return db.Query(sql)
}

//--------File--------
func CreateTableWithFile() {
	sql := `create table if not exists file(
        id int(4) primary key auto_increment not null,
        filepath varchar(255),
        filename varchar(64),
        status int(4),
        createtime int(10)
        );`
	ModifyDB(sql)
}
