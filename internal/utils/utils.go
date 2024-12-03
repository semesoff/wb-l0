package utils

import (
	"encoding/json"
	"fmt"
	"wb-l0/internal/models/order"
)

func EncodeMessage(msg []byte) (order.Order, error) {
	var orderData order.Order
	err := json.Unmarshal(msg, &orderData)
	if err != nil {
		return order.Order{}, fmt.Errorf("error while unmarshalling message: %v", err)
	}
	return orderData, nil
}
