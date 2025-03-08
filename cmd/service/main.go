package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	metrics_lib "github.com/QR-authentication/metrics-lib"
	qrproto "github.com/QR-authentication/qr-proto/qr-proto"

	"github.com/QR-authentication/qr-service/internal/config"
	"github.com/QR-authentication/qr-service/internal/infra"
	"github.com/QR-authentication/qr-service/internal/repository/postgres"
	"github.com/QR-authentication/qr-service/internal/service"
)

func main() {
	cfg := config.MustLoad()

	DBRepo := postgres.New(cfg)
	defer DBRepo.Close()

	metrics, err := metrics_lib.New(cfg.Metrics.Host, cfg.Metrics.Port, cfg.Service.Name, cfg.Platform.Env)
	if err != nil {
		log.Fatal("failed to create metrics object: ", err)
	}
	defer metrics.Disconnect()

	QRService := service.New(DBRepo, cfg.Security.SigningKey)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.MetricsInterceptor(metrics),
		),
	)

	qrproto.RegisterQRServiceServer(grpcServer, QRService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("failed to start TCP listener: %v", err)
	}

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to start gRPC listener: %v", err)
	}
}
