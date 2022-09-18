package apps

import (
	"encoding/json"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func App1(path string, ch chan string) error {
	var users []UserData
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("error %v", err)
		return err
	}

	if err != nil {
		fmt.Printf("error %v", err)
		return err
	}

	err = json.Unmarshal([]byte(file), &users)
	if err != nil {
		fmt.Printf("fail to unmarshal %v", err)
		return err
	}

	b, err := json.Marshal(&users)
	if err != nil {
		fmt.Printf("fail to marshal %v", err)
		return err
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Printf("fail to connect due to %v", err)
		return err
	}
	defer conn.Close()

	rbch, err := conn.Channel()
	if err != nil {
		fmt.Printf("fail to connect to channel %v", err)
		return err
	}

	defer rbch.Close()

	q, err := rbch.QueueDeclare(
		"exercise-queue", // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)

	err = rbch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})

	if err != nil {
		fmt.Printf("fail to connect due to %v", err)
		return err
	}

	fmt.Printf("Successfully sent \n")
	ch <- "Ok"
	return nil
}
