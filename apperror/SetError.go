package apperror

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func SetError(ctx context.Context, err error) *gqlerror.Error {
	if gqlErr, ok := err.(*gqlerror.Error); ok {
		return gqlErr
	}

	if appErr, ok := err.(*AppError); ok {
		gqlErr := graphql.DefaultErrorPresenter(ctx, err)
		gqlErr.Message = appErr.Message
		gqlErr.Extensions = map[string]any{
			"code": appErr.Code,
		}
		return gqlErr
	}

	gqlErr := graphql.DefaultErrorPresenter(ctx, err)
	gqlErr.Extensions = map[string]any{
		"code": "INTERNAL_SERVER_ERROR",
	}
	return gqlErr
}
