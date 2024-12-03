package orders

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"wb-l0/internal/cache"
	"wb-l0/internal/models/order"
)

func HandlerOrders(w http.ResponseWriter, r *http.Request, redisProvider *cache.RedisProvider) {
	orderUID := mux.Vars(r)["order_uid"]
	var orderData order.Order
	if data, err := (*redisProvider).BytesToModel(orderUID); err == nil {
		orderData = data
	} else {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	renderOrdersPage(w, orderData)
}

func renderOrdersPage(w http.ResponseWriter, data interface{}) {
	tmpl, err := template.ParseFiles("web/templates/order.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Fatalf("Error parsing html-files: %v", err)
		return
	}
	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Fatalf("Error executing template: %v", err)
		return
	}
}
