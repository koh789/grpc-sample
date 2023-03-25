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
		fmt.Println("3: exit")
		fmt.Print("please enter >")
		scanner.Scan()
		switch scanner.Text() {
		case "1":
			Hello(client, scanner)
		case "2":
			HelloServerStream(client, scanner)
		case "3":
			fmt.Println("bye.")
			goto M
		}
	}
M:
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
