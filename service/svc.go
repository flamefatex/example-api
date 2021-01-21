package service

type Svc interface {
	Name() string
	Run() (err error)
}
