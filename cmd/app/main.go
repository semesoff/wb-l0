package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"wb-l0/config"
	"wb-l0/internal/cache"
	"wb-l0/internal/db/db"
	"wb-l0/internal/kafka"
	"wb-l0/internal/router"
)

func main() {
	cfg := config.NewConfigManager()                     // Init config
	database := db.NewDatabase(cfg.GetConfig().Database) // Init database
	redisClient := cache.NewRedis(cfg.GetConfig().Redis) // Init redis

	var dbProvider db.DatabaseProvider = database
	var redisProvider cache.RedisProvider = redisClient

	database.RestoreCache(&redisProvider)                                                        // Restore cache from database
	kafka.NewKafkaManager(&dbProvider, &redisProvider).StartKafkaServices(cfg.GetConfig().Kafka) // Init kafka services (consumer, producer)

	r := mux.NewRouter()
	router.InitRouter(r, &redisProvider) // Init routers

	cfgApp := cfg.GetConfig().App
	serverAddress := fmt.Sprintf("%s:%s", cfgApp.Host, cfgApp.Port)
	log.Println("Server is starting on", serverAddress)
	if ok := http.ListenAndServe(serverAddress, r); ok != nil {
		log.Fatalf("Error: %v", ok)
	}
}
