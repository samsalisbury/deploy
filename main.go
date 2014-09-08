package main

import (
	"github.com/opentable/hat"
	"github.com/opentable/ot-go-lib/logging"
	"github.com/opentable/ot-go-lib/service"
)

var (
	log = logging.StandardConfig("deploy").StartupLog(0)
	svc service.Service
)

func main() {
	s, err := hat.NewServer(Root{})
	if err != nil {
		log.Fatal(err)
		return
	}
	svc, err := service.NewHTTPServiceFromEnv("deploy", s.ServeHTTP)
	if err != nil {
		log.Fatal(err)
		return
	}
	svc.Start()
}
