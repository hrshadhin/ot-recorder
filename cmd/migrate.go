package cmd

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"ot-recorder/infrastructure/config"
	"strconv"

	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	DBPgsql = "postgres"
	DBMysql = "mysql"
)

//go:generate cp -r ./../infrastructure/db/migrations ./migrations
//go:embed migrations/*
var migrationsFS embed.FS

//nolint:gochecknoglobals
var (
	migrationDirMap = map[string]string{"sqlite": "sqlite", "mysql": "mysql", "postgres": "pgsql"}
	migrateCmd      = &cobra.Command{
		Use:              "migrate",
		Short:            "migrate database",
		Long:             `migrate database like(postgres, mysql, sqlite)`,
		TraverseChildren: true,
	}
)

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(migrateCmd)

	// add subcommands
	migrateCmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "migration up",
		Long:  `migration up`,
		Run: func(cmd *cobra.Command, args []string) {
			handleSubCommand("up", args)
		},
	})

	migrateCmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "migration down",
		Long:  `migration down`,
		Run: func(cmd *cobra.Command, args []string) {
			handleSubCommand("down", args)
		},
	})
}

func handleSubCommand(name string, args []string) {
	var step = 0

	if len(args) > 0 {
		mStep, err := strconv.Atoi(args[0])
		if err != nil {
			logrus.Printf("%q migration step should be a number\n", args[0])
			os.Exit(1)
		}

		step = mStep
	}

	if err := migrateDatabase(name, step); err != nil {
		logrus.Printf("%v", err)
		os.Exit(1)
	}
}

func migrateDatabase(state string, step int) error {
	driverName, dbURL := getDBDriverAndURL()

	db, err := sql.Open(driverName, dbURL)
	if err != nil {
		return err
	}

	defer db.Close()

	instance, err := getDBInstance(db)
	if err != nil {
		return err
	}

	schemaPath := fmt.Sprintf("migrations/%s", migrationDirMap[config.Get().Database.Type])

	sDriver, err := iofs.New(migrationsFS, schemaPath)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", sDriver, driverName, instance)
	if err != nil {
		return err
	}

	var mErr error

	switch state {
	case "up":
		if step > 0 {
			mErr = m.Steps(step)
		} else {
			mErr = m.Up()
		}
	case "down":
		if step > 0 {
			mErr = m.Steps(-step)
		} else {
			mErr = m.Down()
		}
	}

	if errors.Is(mErr, migrate.ErrNoChange) {
		logrus.Info("no change")
		return nil
	}

	return mErr
}

func getDBDriverAndURL() (driver, dbURL string) {
	cfg := config.Get()
	driver = "sqlite3"
	dbURL = fmt.Sprintf("%s/%s.db", cfg.App.DataPath, cfg.Database.Name)

	switch cfg.Database.Type {
	case DBPgsql:
		driver = "pgx"
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
			cfg.Database.SslMode,
		)
	case DBMysql:
		driver = "mysql"
		dbURL = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?multiStatements=true",
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		)
	}

	return driver, dbURL
}

func getDBInstance(db *sql.DB) (instance database.Driver, err error) {
	switch config.Get().Database.Type {
	case DBPgsql:
		instance, err = pgx.WithInstance(db, &pgx.Config{})
		if err != nil {
			return nil, err
		}
	case DBMysql:
		instance, err = mysql.WithInstance(db, &mysql.Config{})
		if err != nil {
			return nil, err
		}
	default:
		instance, err = sqlite3.WithInstance(db, &sqlite3.Config{})
		if err != nil {
			return nil, err
		}
	}

	return instance, err
}
