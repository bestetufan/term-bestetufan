package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bestetufan/beste-store/config"
	"github.com/bestetufan/beste-store/internal/api/router"
	"github.com/bestetufan/beste-store/internal/domain/migration"
	"github.com/bestetufan/beste-store/pkg/database_handler"
	"github.com/bestetufan/beste-store/pkg/httpserver"
	"github.com/bestetufan/beste-store/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.LogLevel)

	// DB
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBDatabaseName)
	db, err := database_handler.CreateDBConnection(connectionString)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - database_handler.CreateDBConnection: %w", err))
	}

	// Migrate & seed
	if cfg.AutoMigrate {
		err = migration.Execute(db)
		if err != nil {
			l.Fatal(fmt.Errorf("app - Run - migration.Execute: %w", err))
		}
	}

	// GIN & router
	gin.SetMode(cfg.GINMode)
	handler := gin.New()
	router.NewRouter(handler, l, cfg, db)

	// HTTP server
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTPPort))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
