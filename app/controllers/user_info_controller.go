package controllers

import (
	"my_project/app/state"
	ctx "my_project/pkg/context"
	"my_project/templates"
)

func UserInfo(c *ctx.WebCtx, _ *state.AppState) error {
	return c.Render(templates.UserInfoPage())
}
