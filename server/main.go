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

// Сервер
type greetingServer struct {
	pb.GreetServiceServer
}

// SayGreeting Серверный метод куда приходят обращения
func (s *greetingServer) SayGreeting(ctx context.Context, req *pb.NoPram) (*pb.GreetingResponse, error) {
	return &pb.GreetingResponse{
		Message: "Greeting",
	}, nil
}

func (s *greetingServer) SayGreetingServerStreaming(req *pb.NamesList, stream pb.GreetService_SayGreetingServerStreamingServer) error {
	log.Printf("Получен запрос с именами: %v", req.Names)

	for _, name := range req.Names {
		res := &pb.GreetingResponse{Message: "Greeting " + name}
		if err := stream.Send(res); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

func (s *greetingServer) SayGreetingClientStreaming(stream pb.GreetService_SayGreetingClientStreamingServer) error {
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

func (s *greetingServer) SayGreetingBidirectionalStreaming(stream pb.GreetService_SayGreetingBidirectionalStreamingServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("Получен запрос с именем: %v", req.Name)
		res := &pb.GreetingResponse{
			Message: "Greeting" + req.Name,
		}
		log.Printf("Отправка сообщения с именем: %v", res.Message)
		if err := stream.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port) //Запуск прослушки порта
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreetServiceServer(grpcServer, &greetingServer{}) //Регистрация сервиса сервера, сгенерированного командой
	log.Printf("Сервер стартовал: %v", lis.Addr())
	err = grpcServer.Serve(lis) //Запуск обслуживания по порту
	if err != nil {
		log.Fatalf("Ошибка старта сервера: %v", err)
	}
}
