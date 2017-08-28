package services

import (
	"log"

	"anthserver/token"

	"github.com/emicklei/go-restful"
	"net/http"
)

type AuthService struct {
	name string
}

func NewAuthService() (service AuthService) {
	service = AuthService{
		name: "AuthService",
	}
	return service
}

func (a AuthService) Name() string {
	return a.name
}

func (a AuthService) Register(webService *restful.WebService) {
	webService.Route(webService.GET("/token").To(getToken))
	log.Println("hello service regist subpath /token")
}

func getToken(req *restful.Request, resp *restful.Response) {
	log.Printf("token function req [%v]\n", req)

	scopes := token.ParseScopes(req.Request.URL)
	access := token.GetResourceActions(scopes)
	result, err := token.MakeToken("admin", "registry", access)
	if err != nil {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp.WriteAsJson(result)
}
