package validation

import (
	"context"
	"time"
)

type Argument interface {
	setUp(ctx *executionContext)
}

func Nil(isNil bool, constraints ...NilConstraint) ValidatorArgument {
	return NewArgument(validateNil(isNil, constraints))
}

func NilProperty(name string, isNil bool, constraints ...NilConstraint) ValidatorArgument {
	return NewArgument(validateNil(isNil, constraints)).At(PropertyName(name))
}

func Bool(value bool, constraints ...BoolConstraint) ValidatorArgument {
	return NewArgument(validateBool(&value, constraints))
}

func BoolProperty(name string, value bool, constraints ...BoolConstraint) ValidatorArgument {
	return NewArgument(validateBool(&value, constraints)).At(PropertyName(name))
}

func NilBool(value *bool, constraints ...BoolConstraint) ValidatorArgument {
	return NewArgument(validateBool(value, constraints))
}

func NilBoolProperty(name string, value *bool, constraints ...BoolConstraint) ValidatorArgument {
	return NewArgument(validateBool(value, constraints)).At(PropertyName(name))
}

func Number[T Numeric](value T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateNumber(&value, constraints))
}

func NumberProperty[T Numeric](
	name string,
	value T,
	constraints ...NumberConstraint[T],
) ValidatorArgument {
	return NewArgument(validateNumber(&value, constraints)).At(PropertyName(name))
}

func NilNumber[T Numeric](value *T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateNumber(value, constraints))
}

func NilNumberProperty[T Numeric](
	name string,
	value *T,
	constraints ...NumberConstraint[T],
) ValidatorArgument {
	return NewArgument(validateNumber(value, constraints)).At(PropertyName(name))
}

func String(value string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateString(&value, constraints))
}

func StringProperty(name, value string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateString(&value, constraints)).At(PropertyName(name))
}

func NilString(value *string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateString(value, constraints))
}

func NilStringProperty(
	name string,
	value *string,
	constraints ...StringConstraint,
) ValidatorArgument {
	return NewArgument(validateString(value, constraints)).At(PropertyName(name))
}

func Countable(count int, constraints ...CountableConstraint) ValidatorArgument {
	return NewArgument(validateCountable(count, constraints))
}

func CountableProperty(
	name string,
	count int,
	constraints ...CountableConstraint,
) ValidatorArgument {
	return NewArgument(validateCountable(count, constraints)).At(PropertyName(name))
}

func Time(value time.Time, constraints ...TimeConstraint) ValidatorArgument {
	return NewArgument(validateTime(&value, constraints))
}

func TimeProperty(name string, value time.Time, constraints ...TimeConstraint) ValidatorArgument {
	return NewArgument(validateTime(&value, constraints)).At(PropertyName(name))
}

func NilTime(value *time.Time, constraints ...TimeConstraint) ValidatorArgument {
	return NewArgument(validateTime(value, constraints))
}

func NilTimeProperty(
	name string,
	value *time.Time,
	constraints ...TimeConstraint,
) ValidatorArgument {
	return NewArgument(validateTime(value, constraints)).At(PropertyName(name))
}

func Valid(value Validatable) ValidatorArgument {
	return NewArgument(validateIt(value))
}

func ValidProperty(name string, value Validatable) ValidatorArgument {
	return NewArgument(validateIt(value)).At(PropertyName(name))
}

func ValidSlice[T Validatable](values []T) ValidatorArgument {
	return NewArgument(validateSlice(values))
}

func ValidSliceProperty[T Validatable](name string, values []T) ValidatorArgument {
	return NewArgument(validateSlice(values)).At(PropertyName(name))
}

func ValidMap[T Validatable](values map[string]T) ValidatorArgument {
	return NewArgument(validateMap(values))
}

func ValidMapProperty[T Validatable](name string, values map[string]T) ValidatorArgument {
	return NewArgument(validateMap(values)).At(PropertyName(name))
}

func Comparable[T comparable](value T, constraints ...ComparableConstraint[T]) ValidatorArgument {
	return NewArgument(validateComparable(&value, constraints))
}

func ComparableProperty[T comparable](
	name string,
	value T,
	constraints ...ComparableConstraint[T],
) ValidatorArgument {
	return NewArgument(validateComparable(&value, constraints)).At(PropertyName(name))
}

func NilComparable[T comparable](
	value *T,
	constraints ...ComparableConstraint[T],
) ValidatorArgument {
	return NewArgument(validateComparable(value, constraints))
}

func NilComparableProperty[T comparable](
	name string,
	value *T,
	constraints ...ComparableConstraint[T],
) ValidatorArgument {
	return NewArgument(validateComparable(value, constraints)).At(PropertyName(name))
}

func Comparables[T comparable](
	values []T,
	constraints ...ComparablesConstraint[T],
) ValidatorArgument {
	return NewArgument(validateComparables(values, constraints))
}

func ComparablesProperty[T comparable](
	name string,
	values []T,
	constraints ...ComparablesConstraint[T],
) ValidatorArgument {
	return NewArgument(validateComparables(values, constraints)).At(PropertyName(name))
}

func EachString(values []string, constraints ...StringConstraint) ValidatorArgument {
	return NewArgument(validateEachString(values, constraints))
}

func EachStringProperty(
	name string,
	values []string,
	constraints ...StringConstraint,
) ValidatorArgument {
	return NewArgument(validateEachString(values, constraints)).At(PropertyName(name))
}

func EachNumber[T Numeric](values []T, constraints ...NumberConstraint[T]) ValidatorArgument {
	return NewArgument(validateEachNumber(values, constraints))
}

func EachNumberProperty[T Numeric](
	name string,
	values []T,
	constraints ...NumberConstraint[T],
) ValidatorArgument {
	return NewArgument(validateEachNumber(values, constraints)).At(PropertyName(name))
}

func EachComparable[T comparable](
	values []T,
	constraints ...ComparableConstraint[T],
) ValidatorArgument {
	return NewArgument(validateEachComparable(values, constraints))
}

func EachComparableProperty[T comparable](
	name string,
	values []T,
	constraints ...ComparableConstraint[T],
) ValidatorArgument {
	return NewArgument(validateEachComparable(values, constraints)).At(PropertyName(name))
}

func CheckNoViolations(err error) ValidatorArgument {
	return NewArgument(
		func(ctx context.Context, validator *Validator) (*ViolationListError, error) {
			return unwrapViolationList(err)
		},
	)
}

func Check(isValid bool) Checker {
	return Checker{
		isValid:         isValid,
		err:             ErrNotValid,
		messageTemplate: ErrNotValid.Message(),
	}
}

func CheckProperty(name string, isValid bool) Checker {
	return Check(isValid).At(PropertyName(name))
}

type ValidateFunc func(ctx context.Context, validator *Validator) (*ViolationListError, error)

func NewArgument(validate ValidateFunc) ValidatorArgument {
	return ValidatorArgument{validate: validate}
}

func This[T any](v T, constraints ...Constraint[T]) ValidatorArgument {
	return NewArgument(
		func(ctx context.Context, validator *Validator) (*ViolationListError, error) {
			violations := NewViolationList()

			for _, constraint := range constraints {
				err := violations.AppendFromError(constraint.Validate(ctx, validator, v))
				if err != nil {
					return nil, err
				}
			}

			return violations, nil
		},
	)
}

type ValidatorArgument struct {
	validate  ValidateFunc
	path      []PropertyPathElement
	isIgnored bool
}

func (arg ValidatorArgument) At(path ...PropertyPathElement) ValidatorArgument {
	arg.path = append(arg.path, path...)
	return arg
}

func (arg ValidatorArgument) When(condition bool) ValidatorArgument {
	arg.isIgnored = !condition
	return arg
}

func (arg ValidatorArgument) setUp(ctx *executionContext) {
	if !arg.isIgnored {
		ctx.addValidation(arg.validate, arg.path...)
	}
}

type Checker struct {
	err               error
	messageTemplate   string
	path              []PropertyPathElement
	groups            []string
	messageParameters TemplateParameterList
	isIgnored         bool
	isValid           bool
}

func (c Checker) At(path ...PropertyPathElement) Checker {
	c.path = append(c.path, path...)
	return c
}

func (c Checker) When(condition bool) Checker {
	c.isIgnored = !condition
	return c
}

func (c Checker) WhenGroups(groups ...string) Checker {
	c.groups = groups
	return c
}

func (c Checker) WithError(err error) Checker {
	c.err = err
	return c
}

func (c Checker) WithMessage(template string, parameters ...TemplateParameter) Checker {
	c.messageTemplate = template
	c.messageParameters = parameters

	return c
}

func (c Checker) setUp(arguments *executionContext) {
	arguments.addValidation(c.validate, c.path...)
}

func (c Checker) validate(ctx context.Context, validator *Validator) (*ViolationListError, error) {
	if c.isValid || c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return &ViolationListError{}, nil
	}

	violation := validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(c.messageParameters...).
		Create()

	return NewViolationList(violation), nil
}
