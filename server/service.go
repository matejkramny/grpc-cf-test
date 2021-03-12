package main

import (
	"fmt"
	testgrpc "grpc_test/grpc"
	"io"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type svcImpl struct {
	testgrpc.TestServiceServer

	mutex sync.Mutex
	conns map[string]chan *testgrpc.ServerEventMessage
}

func (svc *svcImpl) Bidirectional(server testgrpc.TestService_BidirectionalServer) error {
	id := uuid.Must(uuid.NewV4()).String()
	fmt.Printf("Client subscribed\n")

	svc.mutex.Lock()
	c := make(chan *testgrpc.ServerEventMessage)
	svc.conns[id] = c
	svc.mutex.Unlock()

	recv := make(chan *testgrpc.Ping)
	go func() {
		for {
			ping, err := server.Recv()
			if server.Context().Err() != nil || err == io.EOF {
				fmt.Printf("grpc error %v\n", err)
				return
			} else if err != nil && err != io.EOF {
				panic(err)
			}

			recv <- ping
		}
	}()

	defer func() {
		close(recv)
		close(c)

		svc.mutex.Lock()
		delete(svc.conns, id)
		svc.mutex.Unlock()
	}()

	for {
		select {
		case <-server.Context().Done():
			return nil
		case m := <-c:
			server.Send(m)
		case <-recv:
			fmt.Printf("Got ping from %v at %v\n", id, time.Now())
			server.Send(&testgrpc.ServerEventMessage{
				Event: testgrpc.ServerEvent_ServerPing,
			})
		}
	}
}

func (svc *svcImpl) ServerStream(_ *emptypb.Empty, server testgrpc.TestService_ServerStreamServer) error {
	id := uuid.Must(uuid.NewV4()).String()
	fmt.Printf("Client subscribed\n")

	svc.mutex.Lock()
	c := make(chan *testgrpc.ServerEventMessage)
	svc.conns[id] = c
	svc.mutex.Unlock()

	defer func() {
		close(c)

		svc.mutex.Lock()
		delete(svc.conns, id)
		svc.mutex.Unlock()
	}()

	for {
		select {
		case <-server.Context().Done():
			return nil
		case m := <-c:
			server.Send(m)
		}
	}
}

func newServer() *svcImpl {
	return &svcImpl{
		conns: map[string]chan *testgrpc.ServerEventMessage{},
	}
}
