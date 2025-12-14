package validation

import (
	"context"
)

const DefaultGroup = "default"

type Validatable interface {
	Validate(ctx context.Context, validator *Validator) error
}

type ValidatableFunc func(ctx context.Context, validator *Validator) error

func (f ValidatableFunc) Validate(ctx context.Context, validator *Validator) error {
	return f(ctx, validator)
}

func Filter(violations ...error) error {
	list := &ViolationListError{}

	for _, violation := range violations {
		err := list.AppendFromError(violation)
		if err != nil {
			return err
		}
	}

	return list.AsError()
}
