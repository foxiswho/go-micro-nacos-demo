package main

import (
	"context"
	"fmt"
	nacos "github.com/liangzibo/go-plugins-micro-registry-nacos/v2"
	"github.com/liangzibo/go-plugins-micro-registry-nacos/v2/feign"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	helloworld "go-micro-nacos-demo/proto"
	"go-micro-nacos-demo/sdk"
	"io/ioutil"
	"net/http"
)

func main() {
	//创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         NamespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      NacosHost,
			ContextPath: "/nacos",
			Port:        NacosPort,
			//Scheme:      "http",
		},
	}
	// 创建服务发现客户端
	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})

	r := nacos.NewRegistry(func(options *registry.Options) {
		options.Context = context.WithValue(context.Background(), "naming_client", namingClient)
	})

	// 定义服务，可以传入其它可选参数
	service := micro.NewService(
		micro.Name("my.micro.service.client"),
		micro.Registry(r))
	// HTTP 服务  https://github.com/foxiswho/go-micro-echo-demo 下载他并执行
	//普通HTTP
	demoHttpGet(r)
	//封装 HTTP
	handler1, err := sdk.NewDemoController(func(options *feign.Options) {
		options.Registry = r
		options.Service = httpServerName
	}).GetHandler1("")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(handler1)
		fmt.Println(handler1)
		fmt.Println(handler1)
		fmt.Println(handler1)
	}
	// 创建新的客户端
	greeter := helloworld.NewGreeterService(RpcServerName, service.Client())
	// 调用greeter
	rsp, err := greeter.Hello(context.TODO(), &helloworld.HelloRequest{Name: "John"})
	if err != nil {
		fmt.Println(err)
	}
	//获取所有服务
	fmt.Println(registry.ListServices())
	//获取某一个服务
	services, err := registry.GetService(RpcServerName)
	if err != nil {
		fmt.Println(err)
	}

	//监听服务
	watch, err := registry.Watch()

	fmt.Println(services)
	// 打印响应请求
	fmt.Println(rsp.Greeting)
	go service.Run()
	for {
		result, err := watch.Next()
		if len(result.Action) > 0 {
			fmt.Println(result, err)
		}
	}

}

func demoHttpGet(r registry.Registry) {
	addr, err := feign.GetServiceAddr(r, httpServerName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(addr)
	if len(addr) <= 0 {
		fmt.Println("addr is null")
	} else {
		url := "http://" + addr + "/handler1"
		get, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(get)
		bytes, err := ioutil.ReadAll(get.Body)

		if err != nil {
			fmt.Println("ioutil.ReadAll err=", err)
			return
		}
		fmt.Println(string(bytes))
	}
}
