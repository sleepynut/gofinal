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
		log.Fatal("NO Connection", err)
		return
	}
	defer task.DB.Close()

	err = task.CreateCustomer()
	if err != nil {
		log.Fatal("cannot create customer table", err)
		return
	}

	r := task.SetupRouter()
	r.Run(":2009")
}
