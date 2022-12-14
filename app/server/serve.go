package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"ot-recorder/app/model"
	"syscall"
	"time"

	"ot-recorder/app"
	locationDelivery "ot-recorder/app/location/delivery/http"
	locationMysqlRepo "ot-recorder/app/location/repository/mysql"
	locationPgsqlRepo "ot-recorder/app/location/repository/pgsql"
	locationSqliteRepo "ot-recorder/app/location/repository/sqlite"
	locationUseCase "ot-recorder/app/location/usecase"
	systemDelivery "ot-recorder/app/system/delivery/http"
	systemRepo "ot-recorder/app/system/repository"
	systemUseCase "ot-recorder/app/system/usecase"
	"ot-recorder/infrastructure/config"
	"ot-recorder/infrastructure/db"
	"ot-recorder/infrastructure/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const gracefullShutdownTime = 5

type Server struct {
	ServerReady chan bool
}

func (s *Server) Serve() {
	cfg := config.Get().App

	db.Connect()

	if cfg.Env != "test" {
		defer db.Close()
	}

	e := setupAPIServer(cfg)

	go func() {
		printBanner()

		if err := e.Start(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
			logrus.Errorf(err.Error())
		}
	}()

	// server is ready now, so close the channel
	s.ServerReady <- true

	// signal channel to capture system calls
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-sigCh
	logrus.Info("shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), gracefullShutdownTime*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logrus.Fatalf("failed to gracefully shutdown the server: %s", err)
	}
}

func setupAPIServer(cfg config.AppConfig) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.WriteTimeout = cfg.WriteTimeout
	e.Server.IdleTimeout = cfg.IdleTimeout
	contextTimeout := cfg.ContextTimeout

	if err := middlewares.Attach(e); err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}

	dbClient := db.GetClient()
	dbType := config.Get().Database.Type

	// repository
	sysRepo := systemRepo.NewSystemRepository(dbClient)

	var lRepo model.LocationRepository

	switch dbType {
	case "postgres":
		lRepo = locationPgsqlRepo.NewPgsqlLocationRepository(dbClient)
	case "mysql":
		lRepo = locationMysqlRepo.NewMysqlLocationRepository(dbClient)
	default:
		lRepo = locationSqliteRepo.NewSqliteLocationRepository(dbClient)
	}

	// use cases
	sysUseCase := systemUseCase.NewSystemUsecase(sysRepo)
	lUseCase := locationUseCase.NewLocationUsecase(lRepo, contextTimeout)

	// delivery
	systemDelivery.NewSystemHandler(e, sysUseCase)
	locationDelivery.NewUserHandler(e, lUseCase)

	return e
}

func printBanner() {
	log.SetFlags(0)
	log.Println("=>>")
	log.Println("  ___               _____               _          ____")
	log.Println(" / _ \\__      ___ _|_   _| __ __ _  ___| | _____  |  _ \\ ___  ___")
	log.Println("| | | \\ \\ /\\ / / '_ \\| || '__/ _` |/ __| |/ / __| | |_) / _ \\/ __|")
	log.Println("| |_| |\\ V  V /| | | | || | | (_| | (__|   <\\__ \\ |  _ <  __/ (__ _")
	log.Println(" \\___/  \\_/\\_/ |_| |_|_||_|  \\__,_|\\___|_|\\_\\___/ |_| \\_\\___|\\___(_)")
	log.Println("                                                                              ")
	log.Printf("     Version: %-10s Build time: %-10s", app.Version, app.BuildTime)
	log.Println("<<=")
}
