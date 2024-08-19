package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"blue-admin.com/configs"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

var (
	DBConn *gorm.DB
)

func GormLoggerFile() (*os.File, error) {
	//  Gorm logging file
	gormLogFile, gerr := os.OpenFile("gormblue.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if gerr != nil {
		log.Fatalf("error opening file: %v", gerr)
	}
	return gormLogFile, nil
}

func ReturnSession() (*gorm.DB, error) {

	//  setting up database connection based on DB type
	app_env := configs.AppConfig.Get("DB_TYPE")

	//  This is file to output gorm logger on to
	gormlogger, err := GormLoggerFile()
	if err != nil {
		fmt.Printf("unable to create log file: %v\n", err)
		return nil, err
	}

	gormFileLogger := log.Logger{}
	gormFileLogger.SetOutput(gormlogger)
	gormFileLogger.Writer()

	gormLogger := log.New(gormFileLogger.Writer(), "\r\n", log.LstdFlags|log.Ldate|log.Ltime|log.Lshortfile)
	newLogger := logger.New(
		gormLogger, // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			Colorful:                  true,        // Enable color
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			// ParameterizedQueries:      true,        // Don't include params in the SQL log

		},
	)

	var DBSession *gorm.DB

	switch app_env {
	case "postgres":
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  configs.AppConfig.Get("POSTGRES_URI"),
			PreferSimpleProtocol: true, // disables implicit prepared statement usage,

		}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   newLogger,
			SkipDefaultTransaction:                   true,
		})
		if err != nil {
			panic(err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			fmt.Printf("Error during connecting to database: %v\n", err)
			return nil, err
		}
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Second)

		DBSession = db
	case "sqlite":
		//  this is sqlite connection
		db, _ := gorm.Open(sqlite.Open(configs.AppConfig.Get("SQLLITE_URI")), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   newLogger,
			SkipDefaultTransaction:                   true,
		})

		sqlDB, err := db.DB()
		if err != nil {
			fmt.Printf("Error during connecting to database: %v\n", err)
			return nil, err
		}
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Second)
		DBSession = db
	case "mysql":
		db, _ := gorm.Open(mysql.New(mysql.Config{
			DSN:                       configs.AppConfig.Get("MYSQL_URI"), // data source name
			DefaultStringSize:         256,                                // default size for string fields
			DisableDatetimePrecision:  true,                               // disable datetime precision, which not supported before MySQL 5.6
			DontSupportRenameIndex:    true,                               // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			DontSupportRenameColumn:   true,                               //  when rename column, rename column not supported before MySQL 8, MariaDB
			SkipInitializeWithVersion: false,                              // auto configure based on currently MySQL version
		}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   newLogger,
			SkipDefaultTransaction:                   true,
		})

		sqlDB, err := db.DB()
		if err != nil {
			fmt.Printf("Error during connecting to database: %v\n", err)
			return nil, err
		}
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Second)
		DBSession = db
	case "":
		//  this is sqlite connection
		db, _ := gorm.Open(sqlite.Open("goframe-2.db"), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   newLogger,
			SkipDefaultTransaction:                   true,
		})

		sqlDB, err := db.DB()
		if err != nil {
			fmt.Printf("Error during connecting to database: %v\n", err)
			return nil, err
		}
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Second)
		DBSession = db
	default:
		//  this is sqlite connection
		db, _ := gorm.Open(sqlite.Open(configs.AppConfig.Get("SQLITE_URI")), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   newLogger,
			SkipDefaultTransaction:                   true,
		})

		sqlDB, err := db.DB()
		if err != nil {
			fmt.Printf("Error during connecting to database: %v\n", err)
			return nil, err
		}
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Second)
		DBSession = db

	}

	DBSession.Use(tracing.NewPlugin())
	return DBSession, nil

}
