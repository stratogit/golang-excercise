package apps

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/go-redis/redis"
)

func getBasicStatistic(disVal []UserData) {

	var ulData, dlData, count, hUl, hDl int

	for _, value := range disVal {
		ulData += value.Ul_data_KB
		dlData += value.Dl_data_kB

		count++
		if hUl < value.Ul_data_KB {
			hUl = value.Ul_data_KB
		}
		if hDl < value.Dl_data_kB {
			hDl = value.Dl_data_kB
		}

	}

	avUlData := ulData / count
	avUlDataMB := avUlData / 1000

	avDlData := dlData / count
	avDlDataMB := avDlData / 1000

	fmt.Println("\n| Data      | Throughput KB | Thoughput MB |")
	fmt.Println("|-----------|---------------|--------------|")
	fmt.Printf("|Average ul | %v       | %v         |\n", avUlData, avUlDataMB)
	fmt.Printf("|Highest ul | %v      | %v        |\n", hUl, hUl/1000)
	fmt.Printf("|Average dl | %v   | %v     |\n", avDlData, avDlDataMB)
	fmt.Printf("|Highest dl | %v  | %v    |\n", hDl, hDl/1000)
}

func App3(path string, ch chan string) error {

	var fileVal []UserData

	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("error %v", err)
		return err
	}

	if err != nil {
		fmt.Printf("error %v", err)
		return err
	}

	err = json.Unmarshal([]byte(file), &fileVal)
	if err != nil {
		fmt.Printf("fail to unmarshal %v", err)
		return err
	}

	// fmt.Printf("users %v \n", &users)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// testing connectivity
	_, err = client.Ping().Result()
	if err != nil {
		fmt.Printf("fail to ping redis %v", err)
		return err
	}

	disStr, err := client.Get("exercise").Result()
	if err != nil {
		fmt.Printf("fail to obtain data from redis %v", err)
	}
	fmt.Printf("Successfully get values\n")
	var disVal []UserData
	err = json.Unmarshal([]byte(disStr), &disVal)
	if err != nil {
		fmt.Printf("fail to unmarshal %v", err)
	}

	if !reflect.DeepEqual(fileVal, disVal) {
		return fmt.Errorf("values are different")
	}
	fmt.Printf("Successfully validated data\n")

	getBasicStatistic(disVal)
	ch <- "Ok"
	return nil
}
