package it

import (
	"broker/broker"
	"broker/proto"
	"fmt"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"syscall"
	"testing"
)

type e2eTestSuite struct {
	suite.Suite
	client broker.Client
	server *grpc.Server
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}


func (s *e2eTestSuite) SetupSuite() {
	port := ":8080"
	s.server = broker.NewServer()

	go startServer(s.server, port)

	s.client = broker.NewClient(port)
}

func (s *e2eTestSuite) TearDownSuite() {
	p, _ := os.FindProcess(syscall.Getpid())
	p.Signal(syscall.SIGINT)
}

func (s *e2eTestSuite) Test_EndToEnd_Producer_Consumer() {
	consumer, err := s.client.Consumer()
	s.NoError(err)

	producer, err := s.client.Producer()
	s.NoError(err)

	expectedDate := proto.Data{
		Topic: "test",
		Value: "test_value",
	}
	err = producer.Send(&expectedDate)
	s.NoError(err)

	err = consumer.Send(&proto.ClientMessage{
		Intent: proto.Intent_SUBSCRIBE,
		Id:     "test",
	})
	s.NoError(err)

	//first message is empty response to subscription
	_, err = consumer.Recv()
	s.NoError(err)

	//second message comes from topics queue
	msg2, err := consumer.Recv()
	s.NoError(err)
	s.Equal(expectedDate.Topic, msg2.Data.Topic)
	s.Equal(expectedDate.Value, msg2.Data.Value)

}

func startServer(s *grpc.Server, port string) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("starting server at port %s\n", port)
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}