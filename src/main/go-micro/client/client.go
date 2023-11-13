package main

import (
	"context"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-study/src/main/go-micro/proto"
)

func main() {
	etcdReg := etcd.NewRegistry(
		registry.Addrs("localhost:2379"),
	)

	// create a new service
	service := micro.NewService(micro.Registry(etcdReg))

	// parse command line flags
	service.Init()

	// use the generated client stub
	client := proto.NewGreeterService("greeter", service.Client())

	// make request
	rsp, err := client.Hello(context.Background(), &proto.HelloRequest{Name: "John"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp)
}
