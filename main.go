package main

import (
	"context"
	"encoding/json"
	_ "expvar"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rtbe/clean-rest-api/delivery/web"
	"github.com/rtbe/clean-rest-api/domain/usecase"
	"github.com/rtbe/clean-rest-api/internal/config"
	"github.com/rtbe/clean-rest-api/internal/database"
	"github.com/rtbe/clean-rest-api/internal/database/migrate"
	"github.com/rtbe/clean-rest-api/internal/logger"
	"github.com/rtbe/clean-rest-api/repository/auth"
	"github.com/rtbe/clean-rest-api/repository/order"
	orderitem "github.com/rtbe/clean-rest-api/repository/order_item"
	"github.com/rtbe/clean-rest-api/repository/product"
	"github.com/rtbe/clean-rest-api/repository/user"
)

func main() {
	logger := logger.NewZapLogger()

	if err := run(logger, os.Args[1:]); err != nil {
		// os.Exit(1) will be impliciyly executed.
		logger.Log("fatal", err.Error())
	}
}

func run(logger logger.Logger, args []string) error {
	logger.Log("info", "main      : server starting...")

	if err := godotenv.Load(); err != nil {
		return errors.New("error loading .env file")
	}

	// Set up application configuration.
	cfg := config.New()
	cfgJSON, err := json.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, "error marshalling config")
	}
	logger.Log("info", fmt.Sprintf("config    : %s", cfgJSON))

	// Initialize connections to databases.
	postgreConfig := database.PostgreConfig{
		User:     cfg.DbUser,
		Password: cfg.DbPassword,
		Host:     cfg.DbHost,
	}

	logger.Log("info", "db        : establishing connection to PostgreSQL")
	postgreDB, err := database.NewPostgreSQL(postgreConfig)
	if err != nil {
		return err
	}
	logger.Log("info", "db        : connection to PostgreSQL has been established")
	defer func() {
		if err := postgreDB.Close(); err == nil {
			logger.Log("info", "main      : database disconnected")
		}
	}()
	// Run database migrations.
	logger.Log("info", "db        : running migrations")
	err = migrate.Do(postgreDB)
	if err != nil {
		return err
	}

	mongoConfig := database.MongoConfig{
		User:     cfg.AuthDbUser,
		Password: cfg.AuthDbPassword,
		Host:     cfg.AuthDbHost,
		Name:     cfg.AuthDBName,
	}
	logger.Log("info", "auth db   : establishing connection to MongoDB")
	mongoDB, err := database.NewMongo(mongoConfig)
	if err != nil {
		return err
	}
	logger.Log("info", "auth db   : connection to MongoDB has been established")
	defer func() {
		if err := mongoDB.Client().Disconnect(context.Background()); err == nil {
			logger.Log("info", "main      : auth database disconnected")
		}
	}()
	// Initialize application layers
	userRepo := user.NewPostgreRepo(postgreDB, logger)
	userService := usecase.NewUserService(userRepo)

	productRepo := product.NewPostgreRepo(postgreDB, logger)
	productService := usecase.NewProductService(productRepo)

	orderRepo := order.NewPostgreRepo(postgreDB, logger)
	orderService := usecase.NewOrderService(orderRepo)

	orderItemRepo := orderitem.NewPostgreRepo(postgreDB, logger)
	orderItemService := usecase.NewOrderItemService(orderItemRepo)

	authRepo := auth.NewMongoRepo(mongoDB, logger)
	authService := usecase.NewAuthService(authRepo, userService)

	services := usecase.Services{
		User:      userService,
		Product:   productService,
		Order:     orderService,
		OrderItem: orderItemService,
		Auth:      authService,
	}

	//===============================================Init application server========================================
	app := web.NewApp(services, logger)

	// Configure application server.
	appServer := &http.Server{
		Addr:         ":" + cfg.APIPort,
		Handler:      app,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Init server for debugging on default serve mux,
	// it will be available at: /debug/pprof
	pprofServer := http.Server{
		Addr:         ":8081",
		Handler:      http.DefaultServeMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	go func() {
		logger.Log("info", fmt.Sprintf(
			"main      : pprof and documentation server listening on port %s",
			cfg.APIPort,
		))

		if err := pprofServer.ListenAndServe(); err != nil {
			log.Printf("error serving debug/pprof: %v ", err)
			os.Exit(1)
		}
	}()

	// Gracefull shutdown.
	serverErrors := make(chan error)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Log("info", fmt.Sprintf(
			"main      : server listening on port %s",
			cfg.APIPort,
		))

		serverErrors <- appServer.ListenAndServe()
	}()

	select {

	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		logger.Log("info", fmt.Sprintf(
			"main      : %v: start shutdown",
			sig,
		))

		// Give 30 seconds to shut down, then shut down forcefully.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		appServer.SetKeepAlivesEnabled(false)

		if err := appServer.Shutdown(ctx); err != nil {
			appServer.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}
	return nil
}
