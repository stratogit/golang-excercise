package apps

import (
	"encoding/json"
	"fmt"
	"os"

	"task/src/usersTypes"
	"task/src/utils"
)

var NumberOfUsers int

func App1(path string, ch chan string) error {
	var users []usersTypes.UserData
	var c int

	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("error %v", err)
		return err
	}

	err = json.Unmarshal([]byte(file), &users)
	if err != nil {
		fmt.Printf("fail to unmarshal %v", err)
		return err
	}

	for _, user := range users {
		b, err := json.Marshal(&user)
		if err != nil {
			fmt.Printf("fail to marshal %v", err)
			return err
		}

		c++
		queueName := fmt.Sprintf("exercise-queue-%d", c)
		err = utils.RabbitCreateQueue(queueName, b)
		if err != nil {
			fmt.Printf("fail to create queue %v", err)
			return err
		}
	}

	NumberOfUsers = c
	fmt.Printf("Successfully sent \n")
	ch <- "Ok"
	return nil
}
