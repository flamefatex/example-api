package example

import (
	"context"
	"sync"
	"time"

	"github.com/flamefatex/log"
)

var (
	name = "Example"
	// 单例实体
	instance *exampleSvc
	once     sync.Once
)

type ExampleSvc interface {
	Name() string
}

type exampleSvc struct {
}

func NewExampleSvc() *exampleSvc {
	once.Do(func() {
		instance = &exampleSvc{}
	})

	return instance
}
func ExampleSvcInstance() *exampleSvc {
	return instance
}

func (s *exampleSvc) Name() string {
	return name
}

func (s *exampleSvc) Run(ctx context.Context) (err error) {
	round := 1
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			log.Debugf("example round %d", round)
			round++
			<-ticker.C
		}
	}()
	return
}
