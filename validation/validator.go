package validation

import (
	"context"
	"slices"
	"time"
)

type Validator struct {
	propertyPath     *PropertyPath
	violationFactory ViolationFactory
	groups           []string
}

type ValidatorOptions struct {
	violationFactory ViolationFactory
}

func newValidatorOptions() *ValidatorOptions {
	return &ValidatorOptions{}
}

type ValidatorOption func(options *ValidatorOptions) error

func NewValidator(options ...ValidatorOption) (*Validator, error) {
	var err error

	opts := newValidatorOptions()
	for _, setOption := range options {
		err = setOption(opts)
		if err != nil {
			return nil, err
		}
	}

	if opts.violationFactory == nil {
		opts.violationFactory = NewViolationFactory()
	}

	validator := &Validator{
		violationFactory: opts.violationFactory,
	}

	return validator, nil
}

func SetViolationFactory(factory ViolationFactory) ValidatorOption {
	return func(options *ValidatorOptions) error {
		options.violationFactory = factory

		return nil
	}
}

func (validator *Validator) Validate(ctx context.Context, arguments ...Argument) error {
	execContext := &executionContext{}
	for _, argument := range arguments {
		argument.setUp(execContext)
	}

	violations := &ViolationListError{}

	for _, validate := range execContext.validations {
		vs, err := validate(ctx, validator)
		if err != nil {
			return err
		}

		violations.Join(vs)
	}

	return violations.AsError()
}

func (validator *Validator) ValidateBool(
	ctx context.Context,
	value bool,
	constraints ...BoolConstraint,
) error {
	return validator.Validate(ctx, Bool(value, constraints...))
}

func (validator *Validator) ValidateInt(
	ctx context.Context,
	value int,
	constraints ...NumberConstraint[int],
) error {
	return validator.Validate(ctx, Number(value, constraints...))
}

func (validator *Validator) ValidateFloat(
	ctx context.Context,
	value float64,
	constraints ...NumberConstraint[float64],
) error {
	return validator.Validate(ctx, Number(value, constraints...))
}

func (validator *Validator) ValidateString(
	ctx context.Context,
	value string,
	constraints ...StringConstraint,
) error {
	return validator.Validate(ctx, String(value, constraints...))
}

func (validator *Validator) ValidateStrings(
	ctx context.Context,
	values []string,
	constraints ...ComparablesConstraint[string],
) error {
	return validator.Validate(ctx, Comparables(values, constraints...))
}

func (validator *Validator) ValidateCountable(
	ctx context.Context,
	count int,
	constraints ...CountableConstraint,
) error {
	return validator.Validate(ctx, Countable(count, constraints...))
}

func (validator *Validator) ValidateTime(
	ctx context.Context,
	value time.Time,
	constraints ...TimeConstraint,
) error {
	return validator.Validate(ctx, Time(value, constraints...))
}

func (validator *Validator) ValidateEachString(
	ctx context.Context,
	values []string,
	constraints ...StringConstraint,
) error {
	return validator.Validate(ctx, EachString(values, constraints...))
}

func (validator *Validator) ValidateIt(ctx context.Context, validatable Validatable) error {
	return validator.Validate(ctx, Valid(validatable))
}

func (validator *Validator) WithGroups(groups ...string) *Validator {
	v := validator.copy()
	v.groups = groups

	return v
}

func (validator *Validator) IsAppliedForGroups(groups ...string) bool {
	if len(validator.groups) == 0 {
		if len(groups) == 0 {
			return true
		}

		if slices.Contains(groups, DefaultGroup) {
			return true
		}
	}

	for _, g1 := range validator.groups {
		if len(groups) == 0 {
			if g1 == DefaultGroup {
				return true
			}
		}

		if slices.Contains(groups, g1) {
			return true
		}
	}

	return false
}

func (validator *Validator) IsIgnoredForGroups(groups ...string) bool {
	return !validator.IsAppliedForGroups(groups...)
}

func (validator *Validator) CreateConstraintError(
	constraintName,
	description string,
) *ConstraintError {
	return &ConstraintError{
		ConstraintName: constraintName,
		Path:           validator.propertyPath,
		Description:    description,
	}
}

func (validator *Validator) At(path ...PropertyPathElement) *Validator {
	v := validator.copy()
	v.propertyPath = v.propertyPath.With(path...)

	return v
}

func (validator *Validator) AtProperty(name string) *Validator {
	v := validator.copy()
	v.propertyPath = v.propertyPath.WithProperty(name)

	return v
}

func (validator *Validator) AtIndex(index int) *Validator {
	v := validator.copy()
	v.propertyPath = v.propertyPath.WithIndex(index)

	return v
}

func (validator *Validator) CreateViolation(
	ctx context.Context,
	err error,
	message string,
	path ...PropertyPathElement,
) Violation {
	return validator.BuildViolation(ctx, err, message).At(path...).Create()
}

func (validator *Validator) BuildViolation(
	ctx context.Context,
	err error,
	message string,
) *ViolationBuilder {
	b := NewViolationBuilder(validator.violationFactory).BuildViolation(err, message)
	b = b.SetPropertyPath(validator.propertyPath)

	return b
}

func (validator *Validator) BuildViolationList(ctx context.Context) *ViolationListBuilder {
	b := NewViolationListBuilder(validator.violationFactory)
	b = b.SetPropertyPath(validator.propertyPath)

	return b
}

func (validator *Validator) copy() *Validator {
	return &Validator{
		propertyPath:     validator.propertyPath,
		violationFactory: validator.violationFactory,
		groups:           validator.groups,
	}
}
