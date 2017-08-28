package restserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

const (
	PREFIX_URL  = "/service"
	LISTEN_PORT = 8080
)

/**
 * 定义服务接口，实现此接口即可注册到该RestService中
 */
type Service interface {
	Name() string
	Register(webService *restful.WebService)
}

/**
 * 定义AuthServer
 */
type AuthServer struct {
	container *restful.Container
	services  []Service
}

/**
 * 构造一个AuthServer
 */
func New() (server *AuthServer) {
	server = &AuthServer{
		container: restful.NewContainer(),
	}
	server.container.Filter(requestFilter)
	return server
}

/**
 * 添加具体的Rest服务
 */
func (a *AuthServer) AddService(s interface{}) {
	if service, ok := s.(Service); ok {
		log.Printf("add service %s\n", service.Name())
		a.services = append(a.services, service)
	}
}

/**
 * 启动Rest服务
 */
func (a *AuthServer) Run() {
	webService := new(restful.WebService)
	webService.Path(PREFIX_URL)
	webService.Produces(restful.MIME_JSON)

	log.Printf("service size %d\n", len(a.services))

	for i, service := range a.services {
		service.Register(webService)
		log.Printf("register service %d : %s\n", i, service.Name())
	}

	// 加入服务到容器中
	a.container.Add(webService)

	// 启动服务
	log.Printf("start listening on %d\n", LISTEN_PORT)
	server := &http.Server{Addr: fmt.Sprintf(":%d", LISTEN_PORT), Handler: a.container}
	server.ListenAndServe()
}

/**
 * 全局filter，用户诊断访问
 */
func requestFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Printf("request filter %s, %s\n", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}
