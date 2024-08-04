package context

import (
	"github.com/a-h/templ"
)

type WebCtx struct {
	*Base
}

func NewWebContext(base *Base) *WebCtx {
	return &WebCtx{
        Base: base,
	}
}

func (self *WebCtx) Render(component templ.Component) error {
	self.Set("Content-Type", "text/html")
	return component.Render(self.UserContext(), self.Response().BodyWriter())
}
