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

func callSayHelloClientStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Клиентский стримминг начат")
	stream, err := client.SayHelloClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("Нельзя стартовать клиентский стримминг: %v", err)
	}
	for _, name := range names.Names {
		req := &pb.HelloRequest{
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

func callSayHelloBidirectionalStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Двухсторонний стримиинг начат")
	stream, err := client.SayHelloBidirectionalStreaming(context.Background())
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
		req := &pb.HelloRequest{
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

	//callSayHello(client)
	// callSayHelloServerStream(client, names)
	callSayHelloClientStream(client, names)
}
