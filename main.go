package main

import (
	"fmt"
	"log"

	"github.com/sleepynut/gofinal/task"
)

func main() {
	fmt.Println("customer service")
	err := task.OpenConnection()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer task.DB.Close()

	err = task.CreateCustomer()
	if err != nil {
		fmt.Println("cannot create customer table", err)
		return
	}

	r := task.SetupRouter()
	r.Run(":2009")
}
