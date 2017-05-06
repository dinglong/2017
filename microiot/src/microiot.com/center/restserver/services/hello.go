package services

import (
    "github.com/emicklei/go-restful"
    "io"
    "log"
)

type HelloService struct{}

func (HelloService) Name() string {
    return "HelloService"
}

func (helloService HelloService) Register(webService *restful.WebService) {
    webService.Route(webService.GET("/hello").To(hello))
    log.Println("hello service regist subpath /hello")
}

func hello(req *restful.Request, resp *restful.Response) {
    log.Printf("hello function req [%v]\n", req)
    io.WriteString(resp, "microiot")
}
