package middleware

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/Foxtrot-14/FitRang/profile-service/metrics"
)

type GraphQLPrometheus struct{}

func (g *GraphQLPrometheus) ExtensionName() string {
	return "GraphQLPrometheus"
}

func (g *GraphQLPrometheus) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (g *GraphQLPrometheus) InterceptOperation(
	ctx context.Context,
	next graphql.OperationHandler,
) graphql.ResponseHandler {

	opCtx := graphql.GetOperationContext(ctx)
	start := time.Now()

	respHandler := next(ctx)

	return func(ctx context.Context) *graphql.Response {
		resp := respHandler(ctx)

		if opCtx != nil && opCtx.Operation != nil {
			opType := string(opCtx.Operation.Operation)
			opName := opCtx.Operation.Name
			if opName == "" {
				opName = "anonymous"
			}

			metrics.GraphQLOperationDuration.
				WithLabelValues(opName, opType).
				Observe(time.Since(start).Seconds())
		}

		return resp
	}
}
