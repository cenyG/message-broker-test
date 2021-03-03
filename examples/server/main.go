package main

import (
	"broker/broker"
	"broker/examples"
	"time"
)

func main() {
	s := broker.NewServer()
	ex := examples.NewExample()

	go ex.StartServer(s)

	time.Sleep(60 * time.Second)
}