package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"wb-l0/internal/cache"
	"wb-l0/internal/handlers/orders"
)

func InitRouter(r *mux.Router, redisProvider *cache.RedisProvider) {
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	r.HandleFunc("/orders/{order_uid}", func(writer http.ResponseWriter, request *http.Request) {
		orders.HandlerOrders(writer, request, redisProvider)
	}).Methods("GET")
}
