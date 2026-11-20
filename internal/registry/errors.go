package registry

import "errors"

var (
	ErrServiceInstanceAlreadyExists = errors.New("service instance already exists")
	ErrServiceInstanceNotFound      = errors.New("service instance not found")
	ErrServiceInstanceDeleted       = errors.New("service instance is deleted")
	ErrNotLeader                    = errors.New("registry node is not the leader")
	ErrUnavailable                  = errors.New("registry leader unavailable")
)
