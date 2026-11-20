package registry

import "context"

type Worker[T any] struct {
	items   <-chan T
	errs    chan<- error
	process func(context.Context, T) error
}

func New[T any](items <-chan T, errs chan<- error, process func(context.Context, T) error) *Worker[T] {
	return &Worker[T]{
		items:   items,
		errs:    errs,
		process: process,
	}
}

func (w *Worker[T]) Start(ctx context.Context) {
	go w.eventloop(ctx)
}

func (w *Worker[T]) eventloop(ctx context.Context) {
	for {
		select {
		case item, ok := <-w.items:
			if !ok {
				return
			}
			if err := w.process(ctx, item); err != nil {
				select {
				case w.errs <- err:
				case <-ctx.Done():
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
