package sdk

import (
	"fmt"
	"github.com/liangzibo/go-plugins-micro-registry-nacos/v2/feign"
)

type DemoController interface {
	GetHandler1(data interface{}) error
}
type demoController struct {
	opts    feign.Options
}

func NewDemoController(opts ...feign.Option) *demoController {
	d := new(demoController)
	for _, o := range opts {
		o(&d.opts)
	}
	fmt.Println(d.opts.Service)
	return d
}

func (i *demoController) GetHandler1(data interface{}) (string, error) {
	return feign.GetFeign(i.opts, "/handler1")
}
