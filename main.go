package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp:@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"go_service_req", // queue name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var req struct {
				TypeService string `json:"type_service"`
				Data        string `json:"data"`
			}
			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Printf("failed to parse request: %v", err)
				continue
			}

			res := struct {
				TypeService string `json:"type_service"`
				ResData     string `json:"res_data"`
			}{
				TypeService: req.TypeService,
			}

			// if req.TypeService == "1" {

			// 	result := Message()
			// 	if result == true {
			// 		res.ResData = "true"

			// 	} else {
			// 		res.ResData = "false"
			// 	}
			// 	// convert bool to string

			// }

			resBytes, err := json.Marshal(res)
			if err != nil {
				log.Printf("failed to serialize response: %v", err)
				continue
			}

			if err := ch.Publish(
				"",               // exchange
				"go_service_res", // routing key
				false,            // mandatory
				false,            // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId,
					Body:          resBytes,
				},
			); err != nil {
				log.Printf("failed to publish response: %v", err)
				continue
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
