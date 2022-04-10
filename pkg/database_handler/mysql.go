package database_handler

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var dbConn *gorm.DB

// CreateDBConnection -.
func CreateDBConnection(conString string) (*gorm.DB, error) {
	// Close the existing connection if open
	if dbConn != nil {
		CloseDBConnection(dbConn)
	}

	db, err := gorm.Open(mysql.Open(conString), &gorm.Config{
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return dbConn, fmt.Errorf("mysql_handler - CreateDBConnection - gorm.Open: %w", err)
	}

	sqlDB, err := db.DB()
	sqlDB.SetConnMaxIdleTime(time.Minute * 5)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	dbConn = db

	return dbConn, err
}

// GetDBConnection -.
func GetDBConnection() (*gorm.DB, error) {
	sqlDB, err := dbConn.DB()
	if err != nil {
		return dbConn, fmt.Errorf("mysql_handler - GetDatabaseConnection - DB: %w", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		return dbConn, fmt.Errorf("mysql_handler - GetDatabaseConnection - Ping: %w", err)
	}
	return dbConn, nil
}

// CloseDBConnection -.
func CloseDBConnection(conn *gorm.DB) {
	sqlDB, err := conn.DB()
	if err == nil {
		defer sqlDB.Close()
	}
}

// AutoMigrateDB -.
func AutoMigrateDB(dst ...interface{}) error {
	db, err := GetDBConnection()
	if err != nil {
		return fmt.Errorf("mysql_handler - AutoMigrateDB - GetDatabaseConnection: %w", err)
	}
	err = db.AutoMigrate(dst)
	if err != nil {
		return fmt.Errorf("mysql_handler - AutoMigrateDB - AutoMigrate: %w", err)
	}
	return nil
}
