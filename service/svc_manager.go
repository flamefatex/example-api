package service

import (
	"context"
	"fmt"

	"github.com/flamefatex/log"
)

type svcManager struct {
	svcs []Svc
}

func NewSvcManager() *svcManager {
	return &svcManager{
		svcs: make([]Svc, 0),
	}
}

func (m *svcManager) RegisterSvc(svc Svc) {
	m.svcs = append(m.svcs, svc)
}

func (m *svcManager) Run(ctx context.Context) {
	// 运行
	for _, svc := range m.svcs {
		go func(svc Svc) {
			err := svc.Run(ctx)
			if err != nil {
				err = fmt.Errorf("svc:%s, run err:%w", svc.Name(), err)
				log.Error(err)
			}
		}(svc)
	}
}
