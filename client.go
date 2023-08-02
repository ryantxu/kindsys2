package kindsys2

import "context"

// This is the k8s client
type ResourceClient interface {
	Get(ctx context.Context, id StaticMetadata) (Resource, error)
	Create(ctx context.Context, r Resource) (Resource, error)
	Update(ctx context.Context, r Resource) (Resource, error)
	Delete(ctx context.Context, r Resource) error
	// watch... use k8s directly???
}

// This is used to implement a k8s style asynchronous operator
type ResourceOperator[R Resource] interface {
	OnAdd(ctx context.Context, obj *R) error
	OnUpdate(ctx context.Context, oldObj *R, newObj *R) error
	OnDelete(ctx context.Context, obj *R) error
}
