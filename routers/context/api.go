package context

type ApiCtx struct {
	*Base
}

func NewApiContext(base *Base) *ApiCtx {
	return &ApiCtx{
		Base: base,
	}
}
