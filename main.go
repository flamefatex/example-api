package main

import (
	"context"

	"github.com/flamefatex/config"
	"github.com/flamefatex/log"
	"github.com/flamefatex/log/rotation"
	"go.uber.org/zap"
)

var (
	Version     string
	GitCommit   string
	serviceName string      = "example-api" // 服务名称
	ginAddr                 = ":8002"       // gin 端口
	grpcAddr                = ":8083"       // grpc-gateway grpc端口
	httpAddr                = ":8084"       // grpc-gateway http端口
	zlog        *zap.Logger                 // zap logger
)

func main() {
	ctx := context.Background()
	// config
	config.Init(serviceName)
	// log
	initLogger()
	// print service version
	log.Infof("serviceName: %s, version: %s, build: %s", serviceName, Version, GitCommit)
	// 初始化客户端
	initClient()
	// 初始化并启动svc
	initAndRunSvc(ctx)
	// 启动gin
	runGinServer()
	// 启动grpc gateway
	runGrpcGatewayServer()

}
func initLogger() {
	logConfig := &log.Config{
		ServiceName:    serviceName,
		Level:          config.Config().GetString("log.level"),
		EnableConsole:  config.Config().GetBool("log.enable_console"),
		EnableRotation: false,
	}

	// 如果没有配置文件名，那么还是输出到stdout/stderr
	if config.Config().GetString("log.filename") != "" {
		logConfig.EnableRotation = true
		logConfig.RotationConfig = &rotation.RotationConfig{
			MaxBackups: config.Config().GetInt("log.max_backups"),
			Filename:   config.Config().GetString("log.filename"),
			MaxSize:    config.Config().GetInt("log.max_size"), // MB
		}
	}
	log.InitLogger(log.NewZapLogger, logConfig)

	// zh
	logger, err := log.GetZap()
	if err != nil {
		panic(err)
	}
	zlog = logger
}
