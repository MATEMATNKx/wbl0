package stan

import (
	"encoding/json"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
)

func orderIsValid(order *Order) bool {
	validate := validator.New()
	err := validate.Struct(order)
	return err == nil
}

func (s *Service) handleMessage(m *stan.Msg) {
	var order Order
	log.Printf("Received: %s\n", m)
	if err := json.Unmarshal(m.Data, &order); err != nil {
		log.Printf("[1H] err:\n")
		log.Println(err)
		return
	}
	log.Printf("[1] err:\n")
	if !orderIsValid(&order) {
		log.Printf("[2H] err:\n")
		return
	}
	log.Printf("[2] err:\n")
	dataBytes, err := json.Marshal(order)
	if err != nil {
		log.Printf("[3H] err:\n")
		return
	}
	log.Printf("[3] err:\n")
	data := string(dataBytes)
	log.Printf("data [%s]\n\n", data)
	s.orderSvc.Create(order.OrderUID, data)
	//log.Printf("data [%s]\n", data)
}
