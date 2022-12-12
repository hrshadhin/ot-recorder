package db

import (
	"database/sql"
	"fmt"
	"ot-recorder/infrastructure/config"

	_ "github.com/jackc/pgx/v4/stdlib" // load pgsql driver
	"github.com/sirupsen/logrus"
)

// must call once before server boot to Get() the db instance
func connectPG() (err error) {
	if dbc.DB != nil {
		logrus.Info("db already initialized")
		return nil
	}

	cfg := config.Get().Database
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SslMode,
	)

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	if cfg.MaxOpenConn > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConn)
	}

	if cfg.MaxIdleConn > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConn)
	}

	if cfg.MaxLifeTime > 0 {
		db.SetConnMaxLifetime(cfg.MaxLifeTime)
	}

	dbc.DB = db

	return nil
}
