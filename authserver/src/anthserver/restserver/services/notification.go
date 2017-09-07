package services

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

type EventService struct {
	name string
}

func NewEventService() (service EventService) {
	service = EventService{
		name: "EventService",
	}
	return service
}

func (e EventService) Name() string {
	return e.name
}

func (e EventService) Register(webService *restful.WebService) {
	webService.Route(webService.POST("/event").To(revcEvent))
	log.Println("hello service regist subpath /event")
}

func revcEvent(req *restful.Request, resp *restful.Response) {
	log.Printf("event function req [%v]\n", req)

	// time.Sleep(100 * time.Second)

	// print body
	if data, err := ioutil.ReadAll(req.Request.Body); err != nil {
		log.Printf("event function read req body error %v\n", err)
	} else {
		log.Printf("event function read body [%v]\n", string(data))
	}

	resp.WriteHeader(http.StatusOK)
}
