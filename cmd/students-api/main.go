package main

import (
	"codersGyan/crud/internal/config"
	"codersGyan/crud/internal/http/handlers/students"
	"codersGyan/crud/internal/storage/sqlite"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup

	_, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage Initialzed", slog.String("env", cfg.ENV),slog.String("version","1.0.0"))
	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", students.New())

	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("server start", slog.String("address :", cfg.Addr))
	//fmt.Printf("Server Start at %s", cfg.HTTPServer.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start Server")
		}
	}()

	<-done

	slog.Info("shutting down the Server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errs:= server.Shutdown(ctx)

	if errs != nil {
		slog.Error("failed to shutdown ", slog.String("error", err.Error()))
	}

	slog.Info("Server Shutdown successfully")

}
