package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/koh789/grpc-sample/pkg/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:8080"
)

func main() {
	fmt.Println("start gRPC Client.")
	//gRPCサーバーとのコネクションを確立

	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("Connection failed.")
		return
	}
	defer conn.Close()
	client := pb.NewGreetingServiceClient(conn)

	// 標準入力から文字列を受け取るスキャナを用意
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("1: request Hello - Unary RPC")
		fmt.Println("2: request HelloServerStream - Server streaming RPC")
		fmt.Println("3: request HelloClientStream - Client streaming RPC")
		fmt.Println("4: request HelloBiStreams - BiDirectional streaming RPC")
		fmt.Println("5: exit")
		fmt.Print("please enter >")
		scanner.Scan()
		switch scanner.Text() {
		case "1":
			Hello(client, scanner)
		case "2":
			HelloServerStream(client, scanner)
		case "3":
			HelloClientStream(client, scanner)
		case "4":
			HelloBiStreams(client, scanner)
		case "5":
			fmt.Println("bye.")
			goto M
		}
	}
M:
}

func Hello(client pb.GreetingServiceClient, scanner *bufio.Scanner) {
	fmt.Println("Please enter your name.")
	scanner.Scan()
	name := scanner.Text()

	req := &pb.HelloRequest{
		Name: name,
	}
	res, err := client.Hello(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.GetMessage())
	}
}

func HelloServerStream(client pb.GreetingServiceClient, scanner *bufio.Scanner) {
	fmt.Println("Please enter your name.")
	scanner.Scan()
	req := &pb.HelloRequest{Name: scanner.Text()}
	stream, err := client.HellowServerStream(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("all the responses have already received.")
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
}

func HelloClientStream(client pb.GreetingServiceClient, scanner *bufio.Scanner) {
	fmt.Println("Please enter your name.")
	stream, err := client.HelloClientStream(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	count := 5
	fmt.Printf("Please enter %d names.\n", count)
	for i := 0; i < count; i++ {
		scanner.Scan()
		if err := stream.Send(&pb.HelloRequest{Name: scanner.Text()}); err != nil {
			fmt.Println(err)
			return
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.GetMessage())
	}
}

func HelloBiStreams(client pb.GreetingServiceClient, scanner *bufio.Scanner) {
	stream, err := client.HelloBiStreams(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Please enter names. please type quit when exiting. \n")
	var sendEnd, recvEnd bool
	for !(sendEnd && recvEnd) {
		if !sendEnd {
			scanner.Scan()
			input := scanner.Text()
			switch input {
			case "quit":
				sendEnd = true
				if err := stream.CloseSend(); err != nil {
					fmt.Println(err)
				}
			default:
				if err := stream.Send(&pb.HelloRequest{Name: input}); err != nil {
					fmt.Println(err)
					sendEnd = true
				}
			}
		}
		if !recvEnd {
			if res, err := stream.Recv(); err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Println(err)
				}
				recvEnd = true
			} else {
				fmt.Println(res.GetMessage())
			}
		}
	}
}
