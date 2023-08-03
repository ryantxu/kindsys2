package kindsys2

import "context"

type StorageGetRequest struct {
	Group     string
	Kind      string
	Namespace string
	Name      string
}

type Storage interface {
	GetObject(ctx context.Context, id StaticMetadata, version string) ([]byte, error)
	GetObjectMetadata(ctx context.Context, id StaticMetadata) (CommonMetadata, error)

	// :( we don't want this to be a "public" API.. bytes??
	Create(ctx context.Context, obj Resource) ([]byte, error)
	Update(ctx context.Context, obj Resource) ([]byte, error)
	Delete(ctx context.Context, id StaticMetadata) error

	List(ctx context.Context)
}
