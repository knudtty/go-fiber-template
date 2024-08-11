package controllers

import (
	ctx "my_project/pkg/context"
	"my_project/templates"
)

func UserInfo(c *ctx.WebCtx) error {
	return c.Render(templates.UserInfoPage())
}
