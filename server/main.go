package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":8000"
)

type helloServer struct {
	pb.GreetServiceServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
	grpcServer := grpc.NewServer()

}
