package broker

import (
	"broker/proto"
	"context"
	"google.golang.org/grpc"
	"log"
)

type Client interface {
	Consumer() (proto.Broker_ConsumeClient, error)
	Producer() (proto.Broker_ProduceClient, error)
}

type brokerClient struct {
	client proto.BrokerClient
}

func (c *brokerClient) Consumer() (proto.Broker_ConsumeClient, error) {
	return c.client.Consume(context.Background())
}

func (c *brokerClient) Producer() (proto.Broker_ProduceClient, error) {
	return c.client.Produce(context.Background())
}

func NewClient(port string) Client {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	bc := proto.NewBrokerClient(conn)
	return &brokerClient{client: bc}
}

