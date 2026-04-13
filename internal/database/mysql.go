package database

import (
	"database/sql"
	"fmt"

	"github.com/ariefsibuea/freshmart-api/config"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLConnection(conf config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		conf.MysqlUser,
		conf.MysqlPassword,
		conf.MysqlHost,
		conf.MysqlPort,
		conf.MysqlDatabase,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database connection failed: %w", err)
	}

	db.SetMaxOpenConns(conf.MysqlMaxOpenConns)
	db.SetMaxIdleConns(conf.MysqlMaxIdleConns)
	db.SetConnMaxLifetime(conf.MysqlMaxConnLifetime)
	db.SetConnMaxIdleTime(conf.MysqlMaxConnIdleTime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("mysql ping failed: %w", err)
	}

	return db, nil
}
