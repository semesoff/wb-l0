package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"wb-l0/internal/handlers/orders"
)

func InitRouter(r *mux.Router) {
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	r.HandleFunc("/orders/{order_uid}", orders.HandlerOrders).Methods("GET")
}
