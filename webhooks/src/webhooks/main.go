package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/emicklei/go-restful"
)

func main() {
	ws := new(restful.WebService)
	ws.Route(ws.POST("/webhook").To(webhook))
	restful.Add(ws)
	log.Printf("listen to %d\n", 8080)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func webhook(req *restful.Request, resp *restful.Response) {
	log.Printf("in webhook\n")

	// dump http request
	dump, err := httputil.DumpRequest(req.Request, true)
	if err != nil {
		log.Printf("dump request error %v\n", err)
		http.Error(resp, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	log.Printf("%q", dump)

	// print body
	if data, err := ioutil.ReadAll(req.Request.Body); err != nil {
		log.Printf("read req body error %v\n", err)
	} else {
		log.Printf("read body [%v]\n", string(data))
	}

	// time.Sleep(10 * time.Second)
	resp.WriteHeader(http.StatusOK)

	log.Printf("out webhook\n")
}
