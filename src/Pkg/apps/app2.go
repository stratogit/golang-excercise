package apps

import (
	"encoding/json"
	"fmt"
	"task/src/utils"

	"github.com/go-redis/redis"
)

func App2(ch chan string) error {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	for c := 1; c <= NumberOfUsers; c++ {
		queueName := fmt.Sprintf("exercise-queue-%d", c)

		mUser, err := utils.RabbitConsume(queueName)
		if err != nil {
			fmt.Printf("fail to resturn user %v", err)
			return err
		}

		b, err := json.Marshal(mUser)
		if err != nil {
			fmt.Printf("fail to marshal %v", err)
			return err
		}

		_, err = client.Ping().Result()
		if err != nil {
			fmt.Printf("fail to ping redis %v", err)
			return err
		}

		err = client.Set(queueName, b, 0).Err()
		if err != nil {
			fmt.Printf("fail set values on redis %v", err)
			return err
		}
	}

	fmt.Printf("Successfully set values\n")
	ch <- "Ok"
	return nil
}
