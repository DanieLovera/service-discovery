package registry

import "context"

type StateChangeHandler func(StateChange)

type Backend interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error

	Register(ctx context.Context, params RegisterParams) error
	Update(ctx context.Context, params UpdateParams) error
	Deregister(ctx context.Context, id ServiceInstanceID) error
	Lookup(ctx context.Context, id ServiceInstanceID) ([]ServiceInstance, error)
	Expire(ctx context.Context, id ServiceInstanceID) error
	Recover(ctx context.Context, id ServiceInstanceID) error

	RegisterStateChangeHandler(handler StateChangeHandler)
}
