package service

import "context"

type Svc interface {
	Name() string
	Run(ctx context.Context) (err error)
}
