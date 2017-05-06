package main

import (
    "microiot.com/center/restserver"
    "microiot.com/center/restserver/services"
)

func main() {
    restserver.RegisterService(services.HelloService{})
    restserver.Run()
}
