package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	pb "go-rpc/proto"
	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

type helloServer struct {
	pb.GreetServiceServer
}

// SayHello Серверный метод куда приходят обращения
func (s *helloServer) SayHello(ctx context.Context, req *pb.NoPram) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "Hello",
	}, nil
}

func (s *helloServer) SayHelloServerStreaming(req *pb.NamesList, stream pb.GreetService_SayHelloServerStreamingServer) error {
	log.Printf("Получен запрос с именами: %v", req.Names)

	for _, name := range req.Names {
		res := &pb.HelloResponse{Message: "Hello " + name}
		if err := stream.Send(res); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

func (s *helloServer) SayHelloClientStreaming(stream pb.GreetService_SayHelloClientStreamingServer) error {
	var messages []string
	for {
		req, err := stream.Recv()
		if err != io.EOF {
			return stream.SendAndClose(&pb.MessagesList{
				Messages: messages,
			})
		}
		if err != nil {
			return err
		}
		log.Printf("Получен запрос с именами: %v", req.Name)
		messages = append(messages, req.Name)
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreetServiceServer(grpcServer, &helloServer{})
	log.Printf("Сервер стартовал: %v", lis.Addr())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Ошибка старта сервера: %v", err)
	}
}
