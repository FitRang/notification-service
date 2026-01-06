package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/Foxtrot-14/FitRang/profile-service/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func PrometheusUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		start := time.Now()

		resp, err := handler(ctx, req)

		svc, method := splitMethod(info.FullMethod)
		st := status.Code(err).String()

		metrics.GRPCRequestsTotal.
			WithLabelValues(svc, method, st).
			Inc()

		metrics.GRPCRequestDuration.
			WithLabelValues(svc, method).
			Observe(time.Since(start).Seconds())

		return resp, err
	}
}

func splitMethod(full string) (string, string) {
	parts := strings.Split(full, "/")
	if len(parts) != 3 {
		return "unknown", "unknown"
	}
	return parts[1], parts[2]
}
