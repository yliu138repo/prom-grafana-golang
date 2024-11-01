package main

import (
	"context"
	"fmt"
	"goprometheus/internal/config"
	"goprometheus/internal/httpserver"
	"goprometheus/internal/instrumentation"
	"goprometheus/internal/pgdatabase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("Failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan
		cancel()
	}()

	connStr := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.Driver, cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	pddb, err := pgdatabase.NewPostgresDB(ctx, connStr)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer pddb.Close()

	userHandler := httpserver.NewUserHandler(pddb)

	instrumentation.PrometheusInit()
	http.NewServeMux()
	server := httpserver.NewHTTPServer(userHandler, cfg.WebServer.Port)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Println("Cannot start server" + err.Error())
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("Server shutdown failed " + err.Error())
	}

	log.Printf("Server exiting")

}
