package configs

import "fmt"

const fmMysqlConnection = "%v:%v@tcp(%v:%v)/%v?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"

type MySQL struct {
	UserName string `json:"user_name" env:"MYSQL_USER"`
	PassWord string `json:"pass_word" env:"MYSQL_PASS"`
	Host     string `json:"host" env:"MYSQL_HOST"`
	Port     int    `json:"port" env:"MYSQL_PORT"`
	Database string `json:"database" env:"MYSQL_DB"`
}


func DBConnectionString() string {
	return fmt.Sprintf(fmMysqlConnection, app.MySql.UserName, app.MySql.PassWord, app.MySql.Host, app.MySql.Port, app.MySql.Database)
}

func (db *MySQL) isValid() bool {
	return true
}
