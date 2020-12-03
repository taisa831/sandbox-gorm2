package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/taisa831/sandbox-gorm2/model"
	"gorm.io/driver/sqlite"

	//"gorm.io/driver/sqlite"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnWithLogger() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/gorm-v2?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//SkipDefaultTransaction:                   false,
		//NamingStrategy:                           nil,
		//FullSaveAssociations:                     false,
		Logger:                                   newLogger,
		//NowFunc:                                  nil,
		//DryRun:                                   false,
		//PrepareStmt:                              false,
		//DisableAutomaticPing:                     false,
		//DisableForeignKeyConstraintWhenMigrating: false,
		//AllowGlobalUpdate:                        false,
		//ClauseBuilders:                           nil,
		//ConnPool:                                 nil,
		//Dialector:                                nil,
		//Plugins:                                  nil,
	})

	//db, err := gorm.Open(sqlite.Open("sample.db"), &gorm.Config{
	//	Logger: newLogger,
	//})

	//db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
	//	NamingStrategy: schema.NamingStrategy{TablePrefix: "t_", SingularTable: true},
	//	Logger: newLogger,
	//})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.User{}, &model.Post{}, &model.Company{}, &model.CreditCard{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Conn() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("sample.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var user model.User
	var post model.Post
	err = db.AutoMigrate(&user, &post)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close()

	return db, nil
}
