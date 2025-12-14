package validation

import (
	"context"
	"time"
)

type Numeric interface {
	~float32 | ~float64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Constraint[T any] interface {
	Validate(ctx context.Context, validator *Validator, v T) error
}

type NilConstraint interface {
	ValidateNil(ctx context.Context, validator *Validator, isNil bool) error
}

type BoolConstraint interface {
	ValidateBool(ctx context.Context, validator *Validator, value *bool) error
}

type NumberConstraint[T Numeric] interface {
	ValidateNumber(ctx context.Context, validator *Validator, value *T) error
}

type StringConstraint interface {
	ValidateString(ctx context.Context, validator *Validator, value *string) error
}

type ComparableConstraint[T comparable] interface {
	ValidateComparable(ctx context.Context, validator *Validator, value *T) error
}

type ComparablesConstraint[T comparable] interface {
	ValidateComparables(ctx context.Context, validator *Validator, values []T) error
}

type CountableConstraint interface {
	ValidateCountable(ctx context.Context, validator *Validator, count int) error
}

type TimeConstraint interface {
	ValidateTime(ctx context.Context, validator *Validator, value *time.Time) error
}

type StringFuncConstraint struct {
	err               error
	isValid           func(string) bool
	messageTemplate   string
	groups            []string
	messageParameters TemplateParameterList
	isIgnored         bool
}

func OfStringBy(isValid func(string) bool) StringFuncConstraint {
	return StringFuncConstraint{
		isValid:         isValid,
		err:             ErrNotValid,
		messageTemplate: ErrNotValid.Message(),
	}
}

func (c StringFuncConstraint) WithError(err error) StringFuncConstraint {
	c.err = err
	return c
}

func (c StringFuncConstraint) WithMessage(
	template string,
	parameters ...TemplateParameter,
) StringFuncConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters

	return c
}

func (c StringFuncConstraint) When(condition bool) StringFuncConstraint {
	c.isIgnored = !condition
	return c
}

func (c StringFuncConstraint) WhenGroups(groups ...string) StringFuncConstraint {
	c.groups = groups
	return c
}

func (c StringFuncConstraint) ValidateString(
	ctx context.Context,
	validator *Validator,
	value *string,
) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" ||
		c.isValid(*value) {
		return nil
	}

	return validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		WithParameter("{{ value }}", *value).
		Create()
}
