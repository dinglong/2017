package restserver

import (
    "log"
    "net/http"

    "github.com/emicklei/go-restful"
)

var webService *restful.WebService

/**
 * 初始化WebService
 */
func init() {
    webService = new(restful.WebService)

    // 设置webservice的属性: 基本路径，接受的数据类型，输出的数据类型
    webService.Path("/microiot").
        Doc("microiot rest interface").
        Consumes(restful.MIME_XML, restful.MIME_JSON).
        Produces(restful.MIME_JSON, restful.MIME_XML)
}

/**
 * 定义服务接口，实现此接口即可注册到该RestService中
 */
type Service interface {
    Name() string
    Register(webService *restful.WebService)
}

/**
 * 向RestService中注册一个Service的实现者
 */
func RegisterService(service Service) {
    log.Printf("register service %s\n", service.Name())
    service.Register(webService)
}

/**
 * 启动Rest服务
 */
func Run() {
    // 创建容器
    wsContainer := restful.NewContainer()
    wsContainer.Filter(globalFilter)

    // 加入服务到容器中
    wsContainer.Add(webService)

    // 启动服务
    log.Println("start listening on 8080")
    server := &http.Server{Addr: ":8080", Handler: wsContainer}
    server.ListenAndServe()
}

/**
 * 全局filter，用户诊断访问
 */
func globalFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
    log.Printf("global filter %s, %s\n", req.Request.Method, req.Request.URL)
    chain.ProcessFilter(req, resp)
}
