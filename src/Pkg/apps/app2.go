package apps

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	amqp "github.com/rabbitmq/amqp091-go"
)

func App2(ch chan string) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Printf("fail to connect %v", err)
		return err
	}

	defer conn.Close()

	rabch, err := conn.Channel()
	if err != nil {
		fmt.Printf("fail to create channel %v", err)
		return err
	}

	defer rabch.Close()

	q, err := rabch.QueueDeclare(
		"exercise-queue", // name
		false,            // durable
		false,            // delete when usused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		fmt.Printf("fail to declare %v", err)
		return err
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
		return err
	}

	cUser := make(chan []UserData)

	go func() error {
		var ausers []UserData
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

	fmt.Printf("successfully received \n")

	mUser := <-cUser

	b, err := json.Marshal(mUser)
	if err != nil {
		fmt.Printf("fail to marshal %v", err)
		return err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err = client.Ping().Result()
	if err != nil {
		fmt.Printf("fail to ping redis %v", err)
		return err
	}

	err = client.Set("exercise", b, 0).Err()
	if err != nil {
		fmt.Printf("fail to ping redis %v", err)
		return err
	}

	fmt.Printf("Successfully set values\n")
	ch <- "Ok"
	return nil
}
