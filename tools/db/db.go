package db

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Database struct {
	Instance *gorm.DB
	sqlDb    *sql.DB
}

var DB *Database

func Init() {
	DB = &Database{}
	DB.connect()
}

func (object *Database) Close() {
	object.sqlDb.Close()
}

func (object *Database) connect() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		"root",     //object.conf.UserName,
		"password", //object.conf.Password,
		"mysql",    //object.conf.Host,
		"ethdb",    //object.conf.Name,
		"utf8",     //object.conf.Charset,
		true,       //object.conf.ParseTime,
		"Local",    //object.conf.Locate,
	)
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,       // Disable color
			},
		),
	})
	if err != nil {
		panic(err)
	}

	object.sqlDb, _ = gormDB.DB()
	object.sqlDb.SetMaxOpenConns(20)
	object.sqlDb.SetMaxIdleConns(20)
	object.sqlDb.SetConnMaxLifetime(time.Hour)
	object.Instance = gormDB
}
