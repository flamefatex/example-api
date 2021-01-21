package main

import (
	"context"
	"net"
	"net/http"

	"github.com/flamefatex/config"
	v1 "github.com/flamefatex/example-api/handler/v1"
	v2 "github.com/flamefatex/example-api/handler/v2"
	v2_ext "github.com/flamefatex/example-api/handler/v2/external"
	"github.com/flamefatex/log"
	"github.com/flamefatex/log/rotation"
	protos_v2 "github.com/flamefatex/protos/goout/example-api/v2"
	protos_v2_ext "github.com/flamefatex/protos/goout/example-api/v2/external"
	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
	// config
	config.Init(serviceName)
	// log
	initLogger()
	// print service version
	log.Infof("serviceName: %s, version: %s, build: %s", serviceName, Version, GitCommit)
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

// gin
func runGinServer() {
	// handler 初始化
	engine := gin.Default()
	v1Group := engine.Group("/v1")
	v1.NewDownloadHandler(v1Group)

	// gin Run the server
	go func() {
		log.Debugf("start gin server")

		log.Fatal(engine.Run(ginAddr))
	}()
}

// GrpcGateway
type grpcGWRegister func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

func runGrpcGatewayServer() {
	// grpc
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(zlog),
			grpc_validator.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(zlog),
			grpc_validator.UnaryServerInterceptor(),
		)),
	)

	protos_v2.RegisterExampleServiceServer(grpcServer, v2.NewExampleHandler())

	// external
	protos_v2_ext.RegisterExampleServiceServer(grpcServer, v2_ext.NewExampleHandler())

	go func() {
		log.Debugf("start grpc server")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	regs := []grpcGWRegister{
		protos_v2.RegisterExampleServiceHandlerFromEndpoint,

		// external
		protos_v2_ext.RegisterExampleServiceHandlerFromEndpoint,
	}

	log.Debugf("start grpc-gateway")
	if err := run(regs); err != nil {
		log.Fatal(err)
	}
}

func run(regs []grpcGWRegister) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true, EnumsAsInts: true}))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	for _, reg := range regs {
		err := reg(ctx, mux, grpcAddr, opts)
		if err != nil {
			return err
		}
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(httpAddr, mux)
}
