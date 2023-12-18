package etcd

import (
	"context"
	clientV3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"log"
)

func addEndpoint(addr string) {
	cli, err := clientV3.NewFromURL("http://localhost:2379")
	if err != nil {
		log.Fatalln(err)
	}
	manager, err := endpoints.NewManager(cli, "/foo/bar/my-service")
	if err != nil {
		log.Fatalln(err)
	}
	err = manager.AddEndpoint(context.Background(), "/foo/bar/my-service", endpoints.Endpoint{Addr: "localhost:8090"})
	if err != nil {
		log.Fatalln(err)
	}
}
