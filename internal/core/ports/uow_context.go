package ports

import "context"

// uowContextKey is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type uowContextKey struct{}

// CtxWithUoW adds the provided UnitOfWork to the context.
func CtxWithUoW(ctx context.Context, uow UnitOfWork) context.Context {
	if ctx == nil || uow == nil {
		return ctx
	}
	return context.WithValue(ctx, uowContextKey{}, uow)
}

// UoWFromCtx extracts a UnitOfWork from the context if present.
func UoWFromCtx(ctx context.Context) (UnitOfWork, bool) {
	if ctx == nil {
		return nil, false
	}
	if v := ctx.Value(uowContextKey{}); v != nil {
		if uow, ok := v.(UnitOfWork); ok {
			return uow, true
		}
	}
	return nil, false
}
