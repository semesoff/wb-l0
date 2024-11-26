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
	cfg := config.NewConfigManager()
	database := db.NewDatabase(cfg.GetConfig().Database)
	kafka.NewKafkaManager(database).StartKafkaServices(cfg.GetConfig().Kafka)

	r := mux.NewRouter()
	router.InitRouter(r)

	cfgApp := cfg.GetConfig().App
	serverAddress := fmt.Sprintf("%s:%s", cfgApp.Host, cfgApp.Port)
	log.Println("Server is starting on", serverAddress)
	if ok := http.ListenAndServe(serverAddress, r); ok != nil {
		log.Fatalf("Error: %v", ok)
	}
}
