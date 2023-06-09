package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	pb "github.com/koh789/grpc-sample/pkg/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = 8080
)

func main() {
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
		log.Printf("start unary RPC, server streaming RPC, client streaming RPC port: %v", port)
		server.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	q := <-quit
	log.Printf("stopping gRPC unary_server... signal: %v\n", q)
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

func (sv *greetingServiceImpl) HellowServerStream(req *pb.HelloRequest, stream pb.GreetingService_HellowServerStreamServer) error {
	count := 5
	for i := 0; i < count; i++ {
		if err := stream.Send(&pb.HelloResponse{
			Message: fmt.Sprintf("[%d] Hello, %s", i, req.GetName()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}

func (sv *greetingServiceImpl) HelloClientStream(stream pb.GreetingService_HelloClientStreamServer) error {
	nameList := make([]string, 0)
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return stream.SendAndClose(&pb.HelloResponse{Message: fmt.Sprintf("Hello, %v,", nameList)})
		}
		if err != nil {
			return err
		}
		nameList = append(nameList, req.GetName())
	}
}

func (sv *greetingServiceImpl) HelloBiStreams(stream pb.GreetingService_HelloBiStreamsServer) error {
	for {
		// receive request
		req, err := stream.Recv()
		// 受信した結果errがio.EOFならリクエスト終了
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.HelloResponse{Message: fmt.Sprintf("Hello, %v!", req.GetName())}); err != nil {
			return err
		}
	}
}
