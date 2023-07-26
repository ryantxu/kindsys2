package kindsys2

import "context"

// ???
type ResourceAdmissionHook = func(ctx context.Context, body []byte) (*Resource, error)

// This will be used for app-sdk controllers
type ResourceController[R Resource] interface {
	OnAdd(ctx context.Context, obj *R) error
	OnUpdate(ctx context.Context, oldObj *R, newObj *R) error
	OnDelete(ctx context.Context, obj *R) error
}
