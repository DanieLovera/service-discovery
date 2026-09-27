package registry

import "context"

type Backend interface {
	Register(ctx context.Context, params RegisterParams) error
	Update(ctx context.Context, params UpdateParams) error
	Deregister(ctx context.Context, id ServiceInstanceID) error
	Lookup(ctx context.Context, id ServiceInstanceID) ([]ServiceInstance, error)
	Expire(ctx context.Context, id ServiceInstanceID) error
	Recover(ctx context.Context, id ServiceInstanceID) error
}
