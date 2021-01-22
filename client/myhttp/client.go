package myhttp

import (
	"net/http"
	"sync"

	"github.com/flamefatex/config"
)

var (
	// 单例实体
	instance *http.Client
	once     sync.Once
)

func NewClient() {
	once.Do(func() {
		instance = &http.Client{Timeout: config.Config().GetDuration("http.timeout")}
	})
}

func ClientInstance() *http.Client {
	return instance
}
