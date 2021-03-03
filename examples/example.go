package examples

import (
	"broker/broker"
	"broker/proto"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

type Example struct {
	Host string
	Port  string
	Topics []string
	Count  int
}

func NewExample() Example {
	return Example {
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
		Topics: []string{"one", "two", "three"},
		Count:  100,
	}
}

func (e Example) StartProducer(delay time.Duration) {
	client := broker.NewClient(fmt.Sprintf("%s%s", e.Host, e.Port))
	producer, err := client.Producer()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range e.Topics {
		for i := 0 ; i < e.Count; i++ {
			time.Sleep(delay * time.Millisecond)

			msg := fmt.Sprintf("%s-%d", v, i)
			err = producer.Send(&proto.Data{
				Topic: v,
				Value: msg,
			})
			log.Printf("[producer] send: topic = %s, value = %s", v, msg)

			if err != nil {
				log.Printf("[producer] error: %v", err)
			}
		}
	}
}

func (e Example) StartConsumer(consumerId string) {
	client := broker.NewClient(fmt.Sprintf("%s%s", e.Host, e.Port))
	consumer, err := client.Consumer()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range e.Topics {
		consumer.Send(&proto.ClientMessage{
			Intent: proto.Intent_SUBSCRIBE,
			Id:     v,
		})
	}

	numMsgToReceive := e.Count* len(e.Topics) + len(e.Topics)
	for i := 0 ; i < numMsgToReceive; i++ {
		msg, err := consumer.Recv()
		if err != nil {
			log.Printf("[consumer:%s] error: %v", consumerId, err)
		}

		if msg.Data == nil {
			log.Printf("[consumer:%s] subcribe msg: id = %s", consumerId, msg.Id)
		} else {
			log.Printf("[consumer:%s] msg: id = %s, topic = [%s], value = [%s]", consumerId, msg.Id, msg.Data.Topic, msg.Data.Value)
			consumer.Send(&proto.ClientMessage{
				Intent: proto.Intent_CONFIRM,
				Id:     msg.Id,
			})
		}
	}


}

func (e Example) StartServer(s *grpc.Server) {
	l, err := net.Listen("tcp", e.Port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("starting server at port %s\n", e.Port)
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}