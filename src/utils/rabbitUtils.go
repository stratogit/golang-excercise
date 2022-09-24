package utils

import (
	"encoding/json"
	"fmt"
	"task/src/usersTypes"

	amqp "github.com/rabbitmq/amqp091-go"
)

func RabbitCreateQueue(queueName string, b []byte) error {

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
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		fmt.Printf("fail to declare %v", err)
		return err
	}

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
		fmt.Printf("fail to publish %v", err)
		return err
	}

	return nil
}

func RabbitConsume(queueName string) (usersTypes.UserData, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Printf("fail to connect due to %v", err)
		return usersTypes.UserData{}, err
	}
	defer conn.Close()

	rabch, err := conn.Channel()

	if err != nil {
		fmt.Printf("fail to connect to channel %v", err)
		return usersTypes.UserData{}, err
	}
	defer rabch.Close()
	q, err := rabch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		fmt.Printf("fail to declare %v", err)
		return usersTypes.UserData{}, err
	}

	msgs, err := rabch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		fmt.Printf("fail to consume %v", err)
		return usersTypes.UserData{}, err
	}
	cUser := make(chan usersTypes.UserData)

	go func() error {
		var ausers usersTypes.UserData
		for d := range msgs {
			err = json.Unmarshal([]byte(d.Body), &ausers)
			if err != nil {
				fmt.Printf("fail to unmarshal %v", err)
				return err
			}
			cUser <- ausers
		}
		return nil
	}()
	mUser := <-cUser
	return mUser, nil
}
