package main

import (
	"anthserver/restserver"
	"anthserver/restserver/services"
)

func main() {
	server := restserver.New()
	server.AddService(services.NewAuthService())
	server.AddService(services.NewEventService())
	server.Run()
}
