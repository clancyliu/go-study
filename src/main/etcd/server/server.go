package main

import (
	"context"
	"fmt"
	_ "go-study/src/main/etcd/codec"
	pb "go-study/src/proto"
	clientV3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
)

type Hello struct {
	pb.UnimplementedGreeterServer
}

func (s *Hello) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Println()
	log.Printf("request: %v", request)
	md, _ := metadata.FromIncomingContext(ctx)
	md.Append("world", "hello")
	ctx = metadata.NewIncomingContext(ctx, md)

	return &pb.HelloResponse{
		Message: request.GetName(),
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8090")
	if err != nil {
		log.Println(err)
		return
	}

	cli, err := clientV3.NewFromURL("http://localhost:2379")
	if err != nil {
		log.Fatalln(err)
	}
	manager, err := endpoints.NewManager(cli, "foo/bar/my-service")
	if err != nil {
		log.Fatalln(err)
	}
	err = manager.AddEndpoint(context.Background(), "foo/bar/my-service/v1", endpoints.Endpoint{Addr: "localhost:8090"})
	if err != nil {
		log.Fatalln(err)
	}

	server := grpc.NewServer()

	pb.RegisterGreeterServer(server, &Hello{})
	log.Printf("server listening at %v", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Printf("failed to serve: %v", err)
	}
}
