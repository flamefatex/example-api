package main

import (
	"github.com/flamefatex/example-api/service"
	"github.com/flamefatex/example-api/service/example"
)

func initAndRunSvc() {
	// init
	exampleSvc := example.NewExampleSvc()

	// register
	sm := service.NewSvcManager()
	sm.RegisterSvc(exampleSvc)

	// run
	sm.Run()
}
