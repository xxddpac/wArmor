package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

var (
	Database *sqlx.DB
)

func InitMysql(config *Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.Username, config.Password, config.Addr, config.Port, config.Database)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(time.Minute * 30)
	db.SetConnMaxIdleTime(time.Minute * 5)
	if err = db.Ping(); err != nil {
		return err
	}
	Database = db
	return nil
}
