package db

import (
	"database/sql"
	"fmt"
	"ot-recorder/infrastructure/config"

	_ "github.com/mattn/go-sqlite3" // load sqlite driver
	"github.com/sirupsen/logrus"
)

// must call once before server boot to Get() the db instance
func connectSqlite() (err error) {
	if dbc.DB != nil {
		logrus.Info("db already initialized")
		return nil
	}

	cfg := config.Get()

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s.db", cfg.App.DataPath, cfg.Database.Name))
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	if cfg.Database.MaxOpenConn > 0 {
		db.SetMaxOpenConns(cfg.Database.MaxOpenConn)
	}

	if cfg.Database.MaxIdleConn > 0 {
		db.SetMaxIdleConns(cfg.Database.MaxIdleConn)
	}

	if cfg.Database.MaxLifeTime > 0 {
		db.SetConnMaxLifetime(cfg.Database.MaxLifeTime)
	}

	dbc.DB = db

	return nil
}
