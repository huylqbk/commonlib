package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	mg "github.com/rubenv/sql-migrate"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database interface {
	DB() (*sql.DB, error)
	GetDB() *gorm.DB
	Transaction() (tx *gorm.DB, commit func(), close func())
	Migrate() (int, error)
}

func GormLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
}

func NewPG(dsn string) (Database, error) {
	con, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: GormLogger(),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := con.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Duration(2000) * time.Millisecond)
	return &db{con: con, driver: "postgres"}, err
}

func NewMySQL(dsn string) (Database, error) {
	con, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: GormLogger(),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := con.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Duration(2000) * time.Millisecond)
	return &db{con: con, driver: "mysql"}, err
}

func GetPostgresDSN(username, password, host, port, dbName, ssl string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		username, password,
		host, port,
		dbName, ssl,
	)
}

func GetMysqlDSN(username, password, host, port, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		dbName)
}

type db struct {
	driver string
	con    *gorm.DB
}

func (c *db) DB() (*sql.DB, error) {
	return c.con.DB()
}

func (c *db) GetDB() *gorm.DB {
	return c.con
}

func (c *db) Transaction() (tx *gorm.DB, commit func(), close func()) {
	tx = c.con.Begin()
	close = func() {
		tx.Rollback()
	}

	commit = func() {
		tx.Commit()
	}

	return tx, commit, close
}

func (c *db) Migrate() (int, error) {
	db, err := c.DB()
	if err != nil {
		return 0, err
	}
	// move to env
	strPath := "migrations/auto"

	Path, err := filepath.Abs(strPath)
	if err != nil {
		return 0, err
	}

	migrations := &mg.FileMigrationSource{
		Dir: Path,
	}

	n, err := mg.Exec(db, c.driver, migrations, mg.Up)

	if err != nil {
		log.Println("Error with migration", err)
		return 0, err
	}

	return n, nil
}
