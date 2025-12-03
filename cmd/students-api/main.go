package main

import (
	"codersGyan/crud/internal/config"
	"fmt"
	"log"
	"net/http"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup

	//setup router
	router := http.NewServeMux()

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to students API"))
	})

	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	fmt.Println("Server Start")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("failed to start Server")
	}

}
