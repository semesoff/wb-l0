package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"wb-l0/config"
	"wb-l0/internal/db/db"
	"wb-l0/internal/kafka"
	"wb-l0/internal/router"
)

func main() {
	config.InitConfig()
	db.InitDB()
	defer db.CloseDB()
	kafka.NewKafkaManager().StartKafkaServices()

	r := mux.NewRouter()
	router.InitRouter(r)

	cfg := config.GetConfig().App
	serverAddress := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	log.Println("Server is starting on", serverAddress)
	if ok := http.ListenAndServe(serverAddress, r); ok != nil {
		log.Fatalf("Error: %v", ok)
	}
}
