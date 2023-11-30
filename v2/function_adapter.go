package v2

//
//import "context"
//
//type FunctionAdapter interface {
//	RegisterAdapter(id string, fn AnyFunction)
//	Adapt(id string, ctx context.Context, before any, afterFunc func(ctx context.Context, after any, err error)) (after any, err error)
//}
//
//type functionAdapter struct {
//	adapters map[string]AnyFunction
//}
//
//func NewFunctionAdapter() FunctionAdapter {
//	return &functionAdapter{adapters: make(map[string]AnyFunction)}
//}
//
//func (f *functionAdapter) RegisterAdapter(id string, fn AnyFunction) {
//	f.adapters[id] = fn
//}
//
//func (f *functionAdapter) Adapt(id string, ctx context.Context, before any, afterFunc func(ctx context.Context, after any, err error)) (after any, err error) {
//	defer func(a *any, e *error) {
//		if afterFunc != nil {
//			afterFunc(ctx, *a, *e)
//		}
//	}(&after, &err)
//	if adapter, ok := f.adapters[id]; ok {
//		return adapter(ctx, before)
//	}
//	return before, nil
//}
