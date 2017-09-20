package services

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"anthserver/token"

	"github.com/emicklei/go-restful"
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
	// client address
	log.Printf("client address: %s\n", req.Request.RemoteAddr)

	if strings.HasPrefix(req.Request.RemoteAddr, "192.168.1.50") {
		// http.Error(resp, fmt.Sprint("1.50"), http.StatusInternalServerError)
		// return
	}

	// dump http request header
	dump, err := httputil.DumpRequest(req.Request, true)
	if err != nil {
		log.Printf("dump request error %v\n", err)
		http.Error(resp, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	log.Printf("%q", dump)

	// parse http header, read username and password
	auth := req.Request.Header.Get("Authorization")
	if u, p, ok := parseBasicAuth(auth); ok {
		log.Printf("username: %s, password: %s\n", u, p)
	} else {
		log.Printf("parse basic auth failure")
		// http.Error(resp, fmt.Sprint(err), http.StatusInternalServerError)
		// return
	}

	// generate token
	scopes := token.ParseScopes(req.Request.URL)
	access := token.GetResourceActions(scopes)
	result, err := token.MakeToken("admin", "registry", access)
	if err != nil {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp.WriteAsJson(result)
}

func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}
