package constraint

import (
	"context"
	"strconv"

	"line/validation"
)

type CountConstraint struct {
	minErr                       error
	divisibleErr                 error
	exactErr                     error
	maxErr                       error
	divisibleByMessageTemplate   string
	exactMessageTemplate         string
	maxMessageTemplate           string
	minMessageTemplate           string
	maxMessageParameters         validation.TemplateParameterList
	groups                       []string
	minMessageParameters         validation.TemplateParameterList
	exactMessageParameters       validation.TemplateParameterList
	divisibleByMessageParameters validation.TemplateParameterList
	divisibleBy                  int
	max                          int
	min                          int
	checkDivisible               bool
	isIgnored                    bool
	checkMax                     bool
	checkMin                     bool
}

func newCountConstraint() CountConstraint {
	return CountConstraint{
		minErr:                     validation.ErrTooFewElements,
		maxErr:                     validation.ErrTooManyElements,
		exactErr:                   validation.ErrNotExactCount,
		divisibleErr:               validation.ErrNotDivisibleCount,
		minMessageTemplate:         validation.ErrTooFewElements.Message(),
		maxMessageTemplate:         validation.ErrTooManyElements.Message(),
		exactMessageTemplate:       validation.ErrNotExactCount.Message(),
		divisibleByMessageTemplate: validation.ErrNotDivisibleCount.Message(),
	}
}

func newCountComparison(min, max int, checkMin, checkMax bool) CountConstraint {
	c := newCountConstraint()
	c.min = min
	c.max = max
	c.checkMin = checkMin
	c.checkMax = checkMax

	return c
}

func HasMinCount(min int) CountConstraint {
	return newCountComparison(min, 0, true, false)
}

func HasMaxCount(max int) CountConstraint {
	return newCountComparison(0, max, false, true)
}

func HasCountBetween(min, max int) CountConstraint {
	return newCountComparison(min, max, true, true)
}

func HasExactCount(count int) CountConstraint {
	return newCountComparison(count, count, true, true)
}

func HasCountDivisibleBy(divisor int) CountConstraint {
	c := newCountConstraint()
	c.checkDivisible = true
	c.divisibleBy = divisor

	return c
}

func (c CountConstraint) When(condition bool) CountConstraint {
	c.isIgnored = !condition
	return c
}

func (c CountConstraint) WhenGroups(groups ...string) CountConstraint {
	c.groups = groups
	return c
}

func (c CountConstraint) WithMinError(err error) CountConstraint {
	c.minErr = err
	return c
}

func (c CountConstraint) WithMaxError(err error) CountConstraint {
	c.maxErr = err
	return c
}

func (c CountConstraint) WithExactError(err error) CountConstraint {
	c.exactErr = err
	return c
}

func (c CountConstraint) WithDivisibleError(err error) CountConstraint {
	c.divisibleErr = err
	return c
}

func (c CountConstraint) WithMinMessage(
	template string,
	parameters ...validation.TemplateParameter,
) CountConstraint {
	c.minMessageTemplate = template
	c.minMessageParameters = parameters

	return c
}

func (c CountConstraint) WithMaxMessage(
	template string,
	parameters ...validation.TemplateParameter,
) CountConstraint {
	c.maxMessageTemplate = template
	c.maxMessageParameters = parameters

	return c
}

func (c CountConstraint) WithExactMessage(
	template string,
	parameters ...validation.TemplateParameter,
) CountConstraint {
	c.exactMessageTemplate = template
	c.exactMessageParameters = parameters

	return c
}

func (c CountConstraint) WithDivisibleMessage(
	template string,
	parameters ...validation.TemplateParameter,
) CountConstraint {
	c.divisibleByMessageTemplate = template
	c.divisibleByMessageParameters = parameters

	return c
}

func (c CountConstraint) ValidateCountable(
	ctx context.Context,
	validator *validation.Validator,
	count int,
) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil
	}

	if c.checkDivisible {
		if c.divisibleBy <= 0 {
			return validator.CreateConstraintError(
				"CountConstraint",
				"divisibleBy must be greater than zero",
			)
		}

		if count%c.divisibleBy != 0 {
			return c.newNotDivisibleViolation(ctx, validator, count)
		}
	}

	if c.checkMax && count > c.max {
		return c.newViolation(
			ctx,
			validator,
			count,
			c.max,
			c.maxErr,
			c.maxMessageTemplate,
			c.maxMessageParameters,
		)
	}

	if c.checkMin && count < c.min {
		return c.newViolation(
			ctx,
			validator,
			count,
			c.min,
			c.minErr,
			c.minMessageTemplate,
			c.minMessageParameters,
		)
	}

	return nil
}

func (c CountConstraint) newViolation(
	ctx context.Context,
	validator *validation.Validator,
	count, limit int,
	err error,
	template string,
	parameters validation.TemplateParameterList,
) validation.Violation {
	if c.checkMin && c.checkMax && c.min == c.max {
		template = c.exactMessageTemplate
		parameters = c.exactMessageParameters
		err = c.exactErr
	}

	return validator.BuildViolation(ctx, err, template).
		WithParameters(
			parameters.Prepend(
				validation.TemplateParameter{Key: "{{ count }}", Value: strconv.Itoa(count)},
				validation.TemplateParameter{Key: "{{ limit }}", Value: strconv.Itoa(limit)},
			)...,
		).
		Create()
}

func (c CountConstraint) newNotDivisibleViolation(
	ctx context.Context,
	validator *validation.Validator,
	count int,
) validation.Violation {
	return validator.BuildViolation(ctx, c.divisibleErr, c.divisibleByMessageTemplate).
		WithParameters(
			c.divisibleByMessageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ count }}", Value: strconv.Itoa(count)},
				validation.TemplateParameter{
					Key:   "{{ divisibleBy }}",
					Value: strconv.Itoa(c.divisibleBy),
				},
			)...,
		).
		Create()
}
