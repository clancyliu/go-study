package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	proto2 "go-study/proto"
	clientV3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"time"
)

func grpcCall() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := proto2.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &proto2.HelloRequest{Name: "clancy"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

func etcdCall() {
	cli, err := clientV3.NewFromURL("http://localhost:2379")
	if err != nil {
		log.Fatalln(err)
	}
	etcdResolver, err := resolver.NewBuilder(cli)
	conn, err := grpc.Dial("etcd:///foo/bar/my-service/v1", grpc.WithResolvers(etcdResolver),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	client := proto2.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &proto2.HelloRequest{Name: "clancy"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

func invoke() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.CallContentSubtype("json")))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	response := proto2.HelloResponse{}
	err = conn.Invoke(context.Background(), "/proto.Greeter/SayHello", &proto2.HelloRequest{Name: "clancy"}, &response)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(response)
}

func invokeStr() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype("json")))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	reqStr := `{
		"name": "clancy"
	}`
	var respStr string
	err = conn.Invoke(context.Background(), "/proto.Greeter/SayHello", reqStr, &respStr)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(respStr)

}

type requestUri struct {
	ServiceName string `uri:"serviceName" binding:"required"`
	MethodName  string `uri:"methodName" binding:"required"`
}

func ginCall() {
	r := gin.Default()
	r.GET("/:serviceName/:methodName", func(c *gin.Context) {
		var uri requestUri
		if err := c.ShouldBindUri(&uri); err != nil {
			log.Println(err)
		}

		rawData, err := c.GetRawData()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(uri, rawData)

		cli, err := clientV3.NewFromURL("http://localhost:2379")
		if err != nil {
			log.Fatalln(err)
		}
		etcdResolver, err := resolver.NewBuilder(cli)
		conn, err := grpc.Dial("etcd:///foo/bar/my-service/v1", grpc.WithResolvers(etcdResolver),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
			grpc.WithDefaultCallOptions(grpc.CallContentSubtype("json")))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(conn.GetState())
		defer func() {
			conn.Close()
			fmt.Println(conn.GetState())
		}()

		var respBytes []byte
		//methodStr := "/proto.Greeter/SayHello"
		methodStr := "/" + uri.ServiceName + "/" + uri.MethodName

		//	reqStr := `{
		//	"name": "clancy"
		//}`

		header := metadata.Pairs("hello", "world")
		var trailer metadata.MD

		err = conn.Invoke(context.Background(), methodStr, rawData, &respBytes, grpc.Header(&header), grpc.Trailer(&trailer))
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(respBytes))
		fmt.Println(trailer)

		c.Data(http.StatusOK, "application/json; charset=utf-8", respBytes)
	})
	r.Run(":9080") // listen and serve on 0.0.0.0:8080
}

func main() {

	//invokeStr()
	ginCall()
}
