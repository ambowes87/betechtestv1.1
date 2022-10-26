package main

import (
	"flag"

	"github.com/ambowes87/betechtestv1.1/pkg/logger"
	"github.com/ambowes87/betechtestv1.1/pkg/notifications"
	"github.com/ambowes87/betechtestv1.1/pkg/usersvc"
)

func main() {
	logger.Log("BE Tech Test v1.1 - Alex Bowes")

	userSvcPort := flag.Int("usersvcport", 8080, "port to listen to requests on")
	flag.Parse()

	notificationsBroker := notifications.NewBroker()
	notificationsBroker.Subscribe("user") // TODO: currently does nothing except create topic channel

	svc := usersvc.New("localhost", "/user", *userSvcPort, notificationsBroker)

	err := svc.Start()
	logger.Log(err.Error())
}
