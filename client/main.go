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
func callSayGreeting(client pb.GreetServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // Отмена в конце

	res, err := client.SayGreeting(ctx, &pb.NoPram{})
	if err != nil {
		log.Fatalf("Невозможно приветствовать: %v", err)
	}
	log.Println(res.Message)
}

func callSayGreetingServerStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Стартован сервер потока")
	stream, err := client.SayGreetingServerStreaming(context.Background(), names)
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

func callSayGreetingClientStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Клиентский стримминг начат")
	stream, err := client.SayGreetingClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("Нельзя стартовать клиентский стримминг: %v", err)
	}
	for _, name := range names.Names {
		req := &pb.GreetingRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Ошибка в ходе отправки: %v", err)
		}
		log.Printf("Отправлен запрос с именем: %v", name)
		time.Sleep(2 * time.Second)
	}
	res, err := stream.CloseAndRecv()
	log.Println("Стримминг закончен")
	if err != nil {
		log.Fatalf("Ошибка при получении ответа с сервера: %v", err)
	}
	log.Printf("%v", res.Messages)
}

func callSayGreetingBidirectionalStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Двухсторонний стримиинг начат")
	stream, err := client.SayGreetingBidirectionalStreaming(context.Background())
	if err != nil {
		log.Fatalf("Нельзя стартовать двухсторонний стримминг: %v", err)
	}
	waitc := make(chan struct{})

	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Ошибка при стримминге: %v", err)
			}
			log.Println(message)
		}
		close(waitc)
	}()

	for _, name := range names.Names {
		req := &pb.GreetingRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Ошибка при отправке: %v", err)
		}
		time.Sleep(2 * time.Second)
	}

	stream.CloseSend()
	<-waitc
	log.Println("Двухсторонний стримминг закончен")

}

func main() {
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Нет соединения: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreetServiceClient(conn)

	names := &pb.NamesList{
		Names: []string{"Alice", "Bob", "King"},
	}

	//callSayGreeting(client)
	// callSayGreetingServerStream(client, names)
	callSayGreetingClientStream(client, names)
}
