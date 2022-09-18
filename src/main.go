package main

import (
	"fmt"
	"os"
	"task/src/Pkg/apps"
)

func main() {

	input := os.Args[1]
	chapp1 := make(chan string)
	go apps.App1(input, chapp1)
	result := <-chapp1
	if result == "Ok" {
		chapp2 := make(chan string)
		go apps.App2(chapp2)
		result = <-chapp2
		if result == "Ok" {
			chapp3 := make(chan string)
			go apps.App3(input, chapp3)
			result = <-chapp3
			fmt.Println("End with:", result)
		}
	}

}
