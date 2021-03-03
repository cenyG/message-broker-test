package main

import (
	"broker/examples"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	ex := examples.NewExample()

	delay, err := strconv.ParseInt(os.Getenv("DELAY"), 10, 64)
	if err != nil {
		log.Fatal("DELAY env parse error")
	}
	go ex.StartProducer(time.Duration(delay))

	time.Sleep(60 * time.Second)
}