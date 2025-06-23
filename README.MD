# Использование сервиса GRPC

В коде реализованы простые вызовы rpc, потоковая передача с сервера, потоковая передача с клиента и двунаправленная
потоковая передача с клиента на сервер и обратно на клиент

Алгоритм работы с RPC:
1. Настраивается окружение команды `protoc-gen-go (`https://grpc.io/docs/languages/go/quickstart/`)
2. Описывается `greet.proto`. В данном приложении описаны сервисы
```proto
service GreetService {
   //Запрос и ответ
   rpc SayGreeting (NoPram) returns (GreetingResponse);

//Серверная потоковая передача — клиент отправляет запрос и получает поток ответа
rpc SayGreetingServerStreaming (NamesList) returns (stream GreetingResponse);

//Клиентская потоковая передача — клиент отправляет поток запросов и получает ответ
rpc SayGreetingClientStreaming (stream GreetingRequest) returns (MessagesList);

//Двунаправленный — клиент и сервер общаются друг с другом с помощью потока
rpc SayGreetingBidirectionalStreaming (stream GreetingRequest) returns (stream GreetingResponse);
}
```
3. Генерируются сервисы (grpc-код) с помощью команды `protoc`
4. Добавляются функции в сервер и клиент
