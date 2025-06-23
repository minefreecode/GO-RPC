package main

import (
	"context"
	pb "go-rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

const (
	port = ":8080"
)

// Вызов приветствия
func callSayHello(client pb.GreetServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // Отмена в конце

	res, err := client.SayHello(ctx, &pb.NoPram{})
	if err != nil {
		log.Fatalf("Невозможно приветствовать: %v", err)
	}
	log.Println(res.Message)
}

func callSayHelloServerStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Стартован сервер потока")
	stream, err := client.SayHelloServerStreaming(context.Background(), names)
	if err != nil {
		log.Fatalf("Нельзя отправлять имена: %v", err)
	}
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Ошибка в процессе стримминига: %v", err)
		}
		log.Println(message)
	}
	log.Println("Поток закончен")
}

func main() {
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Нет соединения: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreetServiceClient(conn)

	/* names := &pb.NamesList{
		Names: []string{"Alice", "Bob", "King"},
	} */

	callSayHello(client)
}
