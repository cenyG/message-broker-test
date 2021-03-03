package main

import (
	"broker/examples"
	"os"
	"time"
)

func main() {
	ex := examples.NewExample()
	go ex.StartConsumer(os.Getenv("CONSUMER_ID"))

	time.Sleep(60 * time.Second)
}