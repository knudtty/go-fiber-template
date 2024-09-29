package utils

import (
	"context"

	ctx "my_project/pkg/context"
)

func IsAuthenticated(c context.Context) bool {
	myCtx, ok := c.Value("myCtx").(*ctx.WebCtx)
	return ok && myCtx != nil && myCtx.Doer != nil
}
