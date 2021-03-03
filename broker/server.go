package broker

import (
	"broker/broker/util"
	"broker/log"
	"broker/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"io"
	"sync"
	"time"
)

func NewServer() *grpc.Server {
	grpcServer := grpc.NewServer()

	topics := util.ConcurrentMapData{
		Mux:     sync.RWMutex{},
		Map: make(map[string]chan *proto.Data),
	}
	confirmations := util.ConcurrentMapBool{
		Mux:     sync.RWMutex{},
		Map: make(map[string]chan bool),
	}

	srv := server{
		Topics:        &topics,
		Confirmations: &confirmations,
	}
	proto.RegisterBrokerServer(grpcServer, &srv)

	return grpcServer
}

type server struct {
	Topics        *util.ConcurrentMapData
	Confirmations *util.ConcurrentMapBool
}

func (s *server) Produce(srv proto.Broker_ProduceServer) error {
	ctx := srv.Context()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			req, err := srv.Recv()
			if err == io.EOF {
				// return will close stream from server side
				return nil
			}
			if err != nil {
				log.Err("receive error %v", err)
				continue
			}

			topicsMap := s.Topics

			// create new topic if not yet presented
			if _, ok := topicsMap.Read(req.Topic); !ok {
				topicsMap.Write(req.Topic, make(chan *proto.Data))
			}

			topicChan, ok := topicsMap.Read(req.Topic)
			if !ok {
				log.Warn("cant read topic %s", req.Topic)
				continue
			}

			topicChan <- req
		}
	}

}

func (s *server) Consume(srv proto.Broker_ConsumeServer) error {
	ctx := srv.Context()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// receive data from stream
			req, err := srv.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				log.Err("receive error %v", err)
				continue
			}

			if req.Intent == proto.Intent_SUBSCRIBE {
				topicsMap := s.Topics

				if _, ok := topicsMap.Read(req.Id); !ok {
					topicsMap.Write(req.Id, make(chan *proto.Data))
				}

				resp := proto.Message{
					Id:   req.Id,
					Data: nil,
				}
				if err := srv.Send(&resp); err != nil {
					log.Err("send error %v", err)
				}
				go s.listenTopic(srv, req.Id)

			} else if req.Intent == proto.Intent_CONFIRM {
				val, ok := s.Confirmations.Read(req.Id)
				if !ok {
					log.Info("no key [%s] in confirmations list was fond", req.Id)
				}
				val <- true
			}
		}
	}

}

func (s *server) listenTopic(srv proto.Broker_ConsumeServer, topic string) {
	topicChan, ok := s.Topics.Read(topic)
	if !ok {
		log.Warn("cant read topic %s", topic)
	}

	for {
		incomeMessage := <-topicChan
		newMessage := proto.Message{
			Id:   uuid.NewString(),
			Data: incomeMessage,
		}
		if err := srv.Send(&newMessage); err == nil {
			s.Confirmations.Write(newMessage.Id, make(chan bool))
			go s.confirmMessageRecv(newMessage.Id)
		}
	}
}

func (s *server) confirmMessageRecv(msgId string) {
	confirmChan, ok := s.Confirmations.Read(msgId)
	if !ok {
		log.Warn("no confirmation channel for msg: %s", msgId)
	}
	select {
	case <- confirmChan:
		log.Info("msg %s confirmed", msgId)
	case <-time.After(5 * time.Second):
		log.Warn("timeout msg %s handling", msgId)
	}

	s.Confirmations.Delete(msgId)
}
