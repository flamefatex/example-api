package main

import (
	"context"

	"github.com/flamefatex/example-api/service"
	"github.com/flamefatex/example-api/service/example"
)

func initAndRunSvc(ctx context.Context) {
	// init
	exampleSvc := example.NewExampleSvc()

	// register
	sm := service.NewSvcManager()
	sm.RegisterSvc(exampleSvc)

	// run
	sm.Run(ctx)
}
