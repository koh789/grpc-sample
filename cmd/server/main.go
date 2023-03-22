package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	pb "github.com/koh789/grpc-sample/pkg/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// 2. gRPCサーバーを作成
	server := grpc.NewServer()
	pb.RegisterGreetingServiceServer(server, NewGreetingServiceImpl())
	// for grpc curl
	reflection.Register(server)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		server.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	q := <-quit
	log.Printf("stopping gRPC server... signal: %v\n", q)
	server.GracefulStop()
}

type greetingServiceImpl struct {
	pb.UnimplementedGreetingServiceServer
}

func NewGreetingServiceImpl() *greetingServiceImpl {
	return &greetingServiceImpl{}
}

func (sv *greetingServiceImpl) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}
