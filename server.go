package main

import (
	"context"
	helloworld "go-micro-nacos-demo/proto"
	"strconv"

	nacos "github.com/liangzibo/go-plugins-micro-registry-nacos/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
)

const (
	NamespaceId = "9d5d3937-27a6-45a4-b300-e30dc3656a90"
	NacosHost   = "192.168.0.254"
	NacosPort   = 8848
	//Rpc 微服务
	RpcServerName = "my.micro.service"
	//http 微服务
	httpServerName = "go.micro.web.echo"
)

type Helloworld struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Helloworld) Hello(ctx context.Context, req *helloworld.HelloRequest, rsp *helloworld.HelloResponse) error {
	logger.Info("Received Helloworld.Call request")
	return nil
}
func main() {
	//命名空间
	nacos.SetNamespaceId(NamespaceId)
	registry := nacos.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{NacosHost + ":" + strconv.Itoa(NacosPort)}
	})
	service := micro.NewService(
		// Set service name
		micro.Name(RpcServerName),
		micro.Address(":8070"),
		micro.Version("0.0.1"),
		// Set service registry
		micro.Registry(registry),
	)
	helloworld.RegisterGreeterHandler(service.Server(), new(Helloworld))
	service.Run()
}
