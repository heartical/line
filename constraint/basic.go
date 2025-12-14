package constraint

import (
	"context"
	"time"

	"line/validation"
)

type NotBlankConstraint[T comparable] struct {
	blank T
	validation.BaseConstraint
	allowNil bool
}

func IsNotBlank() NotBlankConstraint[string] {
	return IsNotBlankComparable[string]()
}

func IsNotBlankNumber[T validation.Numeric]() NotBlankConstraint[T] {
	return IsNotBlankComparable[T]()
}

func IsNotBlankComparable[T comparable]() NotBlankConstraint[T] {
	return NotBlankConstraint[T]{
		BaseConstraint: validation.BaseConstraint{
			Err:             validation.ErrIsBlank,
			MessageTemplate: validation.ErrIsBlank.Message(),
		},
	}
}

func (c NotBlankConstraint[T]) WithAllowedNil() NotBlankConstraint[T] {
	c.allowNil = true
	return c
}

func (c NotBlankConstraint[T]) ValidateString(
	ctx context.Context,
	validator *validation.Validator,
	value *string,
) error {
	if c.ShouldSkip(validator) {
		return nil
	}

	if c.allowNil && value == nil {
		return nil
	}

	if value != nil && *value != "" {
		return nil
	}

	return c.NewViolation(ctx, validator)
}

func (c NotBlankConstraint[T]) ValidateComparable(
	ctx context.Context,
	validator *validation.Validator,
	value *T,
) error {
	if c.ShouldSkip(validator) {
		return nil
	}

	if c.allowNil && value == nil {
		return nil
	}

	if value != nil && *value != c.blank {
		return nil
	}

	return c.NewViolation(ctx, validator)
}

func (c NotBlankConstraint[T]) ValidateBool(
	ctx context.Context,
	validator *validation.Validator,
	value *bool,
) error {
	if c.ShouldSkip(validator) {
		return nil
	}

	if c.allowNil && value == nil {
		return nil
	}

	if value != nil && *value {
		return nil
	}

	return c.NewViolation(ctx, validator)
}

func (c NotBlankConstraint[T]) ValidateTime(
	ctx context.Context,
	validator *validation.Validator,
	value *time.Time,
) error {
	if c.ShouldSkip(validator) {
		return nil
	}

	if c.allowNil && value == nil {
		return nil
	}

	if value != nil && !value.IsZero() {
		return nil
	}

	return c.NewViolation(ctx, validator)
}

func (c NotBlankConstraint[T]) ValidateCountable(
	ctx context.Context,
	validator *validation.Validator,
	count int,
) error {
	if c.ShouldSkip(validator) || count > 0 {
		return nil
	}

	return c.NewViolation(ctx, validator)
}

type BlankConstraint[T comparable] struct {
	blank T
	validation.BaseConstraint
}

func IsBlank() BlankConstraint[string] {
	return IsBlankComparable[string]()
}

func IsBlankNumber[T validation.Numeric]() BlankConstraint[T] {
	return IsBlankComparable[T]()
}

func IsBlankComparable[T comparable]() BlankConstraint[T] {
	return BlankConstraint[T]{
		BaseConstraint: validation.BaseConstraint{
			Err:             validation.ErrNotBlank,
			MessageTemplate: validation.ErrNotBlank.Message(),
		},
	}
}

func (c BlankConstraint[T]) ValidateString(
	ctx context.Context,
	validator *validation.Validator,
	value *string,
) error {
	if c.ShouldSkip(validator) || value == nil || *value == "" {
		return nil
	}

	return c.NewViolation(ctx, validator)
}

func (c BlankConstraint[T]) ValidateComparable(
	ctx context.Context,
	validator *validation.Validator,
	value *T,
) error {
	if c.ShouldSkip(validator) || value == nil || *value == c.blank {
		return nil
	}

	return c.NewViolation(ctx, validator)
}

func (c BlankConstraint[T]) ValidateBool(
	ctx context.Context,
	validator *validation.Validator,
	value *bool,
) error {
	if c.ShouldSkip(validator) || value == nil || !*value {
		return nil
	}

	return c.NewViolation(ctx, validator)
}

func (c BlankConstraint[T]) ValidateTime(
	ctx context.Context,
	validator *validation.Validator,
	value *time.Time,
) error {
	if c.ShouldSkip(validator) || value == nil || value.IsZero() {
		return nil
	}

	return c.NewViolation(ctx, validator)
}

func (c BlankConstraint[T]) ValidateCountable(
	ctx context.Context,
	validator *validation.Validator,
	count int,
) error {
	if c.ShouldSkip(validator) || count == 0 {
		return nil
	}

	return c.NewViolation(ctx, validator)
}

type NotNilConstraint[T comparable] struct {
	validation.BaseConstraint
}

func IsNotNil() NotNilConstraint[string] {
	return IsNotNilComparable[string]()
}

func IsNotNilNumber[T validation.Numeric]() NotNilConstraint[T] {
	return IsNotNilComparable[T]()
}

func IsNotNilComparable[T comparable]() NotNilConstraint[T] {
	return NotNilConstraint[T]{
		BaseConstraint: validation.BaseConstraint{
			Err:             validation.ErrIsNil,
			MessageTemplate: validation.ErrIsNil.Message(),
		},
	}
}

func (c NotNilConstraint[T]) ValidateNil(
	ctx context.Context,
	validator *validation.Validator,
	isNil bool,
) error {
	if c.ShouldSkip(validator) || !isNil {
		return nil
	}

	return c.NewViolation(ctx, validator)
}

func (c NotNilConstraint[T]) ValidateString(
	ctx context.Context,
	validator *validation.Validator,
	value *string,
) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NotNilConstraint[T]) ValidateComparable(
	ctx context.Context,
	validator *validation.Validator,
	value *T,
) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NotNilConstraint[T]) ValidateBool(
	ctx context.Context,
	validator *validation.Validator,
	value *bool,
) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NotNilConstraint[T]) ValidateTime(
	ctx context.Context,
	validator *validation.Validator,
	value *time.Time,
) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

type NilConstraint[T comparable] struct {
	validation.BaseConstraint
}

func IsNil() NilConstraint[string] {
	return IsNilComparable[string]()
}

func IsNilNumber[T validation.Numeric]() NilConstraint[T] {
	return IsNilComparable[T]()
}

func IsNilComparable[T comparable]() NilConstraint[T] {
	return NilConstraint[T]{
		BaseConstraint: validation.BaseConstraint{
			Err:             validation.ErrNotNil,
			MessageTemplate: validation.ErrNotNil.Message(),
		},
	}
}

func (c NilConstraint[T]) ValidateNil(
	ctx context.Context,
	validator *validation.Validator,
	isNil bool,
) error {
	if c.ShouldSkip(validator) || isNil {
		return nil
	}

	return c.NewViolation(ctx, validator)
}

func (c NilConstraint[T]) ValidateString(
	ctx context.Context,
	validator *validation.Validator,
	value *string,
) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NilConstraint[T]) ValidateComparable(
	ctx context.Context,
	validator *validation.Validator,
	value *T,
) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NilConstraint[T]) ValidateBool(
	ctx context.Context,
	validator *validation.Validator,
	value *bool,
) error {
	return c.ValidateNil(ctx, validator, value == nil)
}

func (c NilConstraint[T]) ValidateTime(
	ctx context.Context,
	validator *validation.Validator,
	value *time.Time,
) error {
	return c.ValidateNil(ctx, validator, value == nil)
}
