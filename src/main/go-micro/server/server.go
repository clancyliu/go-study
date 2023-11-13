package main

import (
	"context"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-study/src/main/go-micro/proto"
	"log"
)

func main() {
	etcdReg := etcd.NewRegistry(
		registry.Addrs("localhost:2379"),
	)

	//newRegistry := registry.NewRegistry(registry.Addrs("localhost:2379"))
	// create service
	service := micro.NewService(
		micro.Name("greeter"),
		micro.Version("latest"),
		micro.Registry(etcdReg),
	)

	// initialise flags
	service.Init()

	// register handler
	if err := proto.RegisterGreeterHandler(service.Server(), &Greeter{}); err != nil {
		log.Fatal(err)
	}

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	fmt.Println("hello world")
	return nil
}
