package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/ambowes87/betechtestv1.1/pkg/logger"
	"github.com/ambowes87/betechtestv1.1/pkg/usersvc"
)

const userEndpoint = "/user"

func main() {
	logger.Log("BE Tech Test v1.1 - Alex Bowes")

	port := flag.Int("port", 8080, "port to listen to requests on")
	flag.Parse()

	http.HandleFunc(userEndpoint, usersvc.HandleRequest)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), nil)
	logger.Log(err.Error())
}
