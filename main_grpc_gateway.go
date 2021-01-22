package main

import (
	"context"
	"net"
	"net/http"

	v2 "github.com/flamefatex/example-api/handler/v2"
	v2_ext "github.com/flamefatex/example-api/handler/v2/external"
	"github.com/flamefatex/example-api/service/example"

	"github.com/flamefatex/log"
	protos_v2 "github.com/flamefatex/protos/goout/example-api/v2"
	protos_v2_ext "github.com/flamefatex/protos/goout/example-api/v2/external"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

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

	protos_v2.RegisterExampleServiceServer(grpcServer, v2.NewExampleHandler(example.ExampleSvcInstance()))

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
