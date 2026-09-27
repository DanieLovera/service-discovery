package registry

import "errors"

var (
	ErrServiceInstanceAlreadyExists = errors.New("service instance already exists")
	ErrServiceInstanceNotFound      = errors.New("service instance not found")
	ErrServiceInstanceDeleted       = errors.New("service instance is deleted")
	ErrServiceNotFound              = errors.New("service not found")
)
