package main

import (
    "microiot.com/center/restserver"
    "microiot.com/center/restserver/services"
)

func main() {
    server := restserver.New()
    server.AddService(services.NewHelloService())
    server.Run()
}
