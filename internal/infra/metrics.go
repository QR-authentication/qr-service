package infra

import (
	"context"
	"strings"
	"time"

	"github.com/QR-authentication/qr-service/internal/config"

	metrics "github.com/QR-authentication/metrics-lib"
	"google.golang.org/grpc"
)

func MetricsInterceptor(metrics *metrics.Metrics) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()
		method := strings.Trim(strings.ReplaceAll(info.FullMethod, "/", "_"), "_")
		metrics.Increment(method)

		ctx = context.WithValue(ctx, config.KeyMetrics, metrics)
		resp, err := handler(ctx, req)

		if err != nil {
			metrics.Increment(method + "_error")
		}

		metrics.Duration(time.Since(startTime).Milliseconds(), method)

		return resp, err
	}
}
