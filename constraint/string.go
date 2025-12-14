package constraint

import (
	"context"
	"regexp"
	"strconv"
	"unicode/utf8"

	"line/predicate"
	"line/validation"
)

type LengthConstraint struct {
	minErr                 error
	exactErr               error
	maxErr                 error
	minMessageTemplate     string
	maxMessageTemplate     string
	exactMessageTemplate   string
	groups                 []string
	minMessageParameters   validation.TemplateParameterList
	maxMessageParameters   validation.TemplateParameterList
	exactMessageParameters validation.TemplateParameterList
	max                    int
	min                    int
	checkMax               bool
	checkMin               bool
	isIgnored              bool
}

func newLengthConstraint(min, max int, checkMin, checkMax bool) LengthConstraint {
	return LengthConstraint{
		min:                  min,
		max:                  max,
		checkMin:             checkMin,
		checkMax:             checkMax,
		minErr:               validation.ErrTooShort,
		maxErr:               validation.ErrTooLong,
		exactErr:             validation.ErrNotExactLength,
		minMessageTemplate:   validation.ErrTooShort.Message(),
		maxMessageTemplate:   validation.ErrTooLong.Message(),
		exactMessageTemplate: validation.ErrNotExactLength.Message(),
	}
}

func HasMinLength(min int) LengthConstraint {
	return newLengthConstraint(min, 0, true, false)
}

func HasMaxLength(max int) LengthConstraint {
	return newLengthConstraint(0, max, false, true)
}

func HasLengthBetween(min, max int) LengthConstraint {
	return newLengthConstraint(min, max, true, true)
}

func HasExactLength(count int) LengthConstraint {
	return newLengthConstraint(count, count, true, true)
}

func (c LengthConstraint) When(condition bool) LengthConstraint {
	c.isIgnored = !condition
	return c
}

func (c LengthConstraint) WhenGroups(groups ...string) LengthConstraint {
	c.groups = groups
	return c
}

func (c LengthConstraint) WithMinError(err error) LengthConstraint {
	c.minErr = err
	return c
}

func (c LengthConstraint) WithMaxError(err error) LengthConstraint {
	c.maxErr = err
	return c
}

func (c LengthConstraint) WithExactError(err error) LengthConstraint {
	c.exactErr = err
	return c
}

func (c LengthConstraint) WithMinMessage(
	template string,
	parameters ...validation.TemplateParameter,
) LengthConstraint {
	c.minMessageTemplate = template
	c.minMessageParameters = parameters

	return c
}

func (c LengthConstraint) WithMaxMessage(
	template string,
	parameters ...validation.TemplateParameter,
) LengthConstraint {
	c.maxMessageTemplate = template
	c.maxMessageParameters = parameters

	return c
}

func (c LengthConstraint) WithExactMessage(
	template string,
	parameters ...validation.TemplateParameter,
) LengthConstraint {
	c.exactMessageTemplate = template
	c.exactMessageParameters = parameters

	return c
}

func (c LengthConstraint) ValidateString(
	ctx context.Context,
	validator *validation.Validator,
	value *string,
) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}

	count := len(*value)
	if !utf8.ValidString(*value) {
		count = utf8.RuneCountInString(*value)
	}

	if c.checkMax && count > c.max {
		return c.newViolation(
			ctx,
			validator,
			count,
			c.max,
			*value,
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
			*value,
			c.minErr,
			c.minMessageTemplate,
			c.minMessageParameters,
		)
	}

	return nil
}

func (c LengthConstraint) newViolation(
	ctx context.Context,
	validator *validation.Validator,
	count int,
	limit int,
	value string,
	err error,
	template string,
	parameters validation.TemplateParameterList,
) validation.Violation {
	if c.checkMin && c.checkMax && c.min == c.max {
		err = c.exactErr
		template = c.exactMessageTemplate
		parameters = c.exactMessageParameters
	}

	return validator.
		BuildViolation(ctx, err, template).
		WithParameters(
			parameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: strconv.Quote(value)},
				validation.TemplateParameter{Key: "{{ length }}", Value: strconv.Itoa(count)},
				validation.TemplateParameter{Key: "{{ limit }}", Value: strconv.Itoa(limit)},
			)...,
		).
		Create()
}

type RegexpConstraint struct {
	err               error
	regex             *regexp.Regexp
	messageTemplate   string
	groups            []string
	messageParameters validation.TemplateParameterList
	isIgnored         bool
	match             bool
}

func Matches(regex *regexp.Regexp) RegexpConstraint {
	return RegexpConstraint{
		regex:           regex,
		match:           true,
		err:             validation.ErrNotValid,
		messageTemplate: validation.ErrNotValid.Message(),
	}
}

func DoesNotMatch(regex *regexp.Regexp) RegexpConstraint {
	return RegexpConstraint{
		regex:           regex,
		match:           false,
		err:             validation.ErrNotValid,
		messageTemplate: validation.ErrNotValid.Message(),
	}
}

func (c RegexpConstraint) WithError(err error) RegexpConstraint {
	c.err = err
	return c
}

func (c RegexpConstraint) WithMessage(
	template string,
	parameters ...validation.TemplateParameter,
) RegexpConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters

	return c
}

func (c RegexpConstraint) When(condition bool) RegexpConstraint {
	c.isIgnored = !condition
	return c
}

func (c RegexpConstraint) WhenGroups(groups ...string) RegexpConstraint {
	c.groups = groups
	return c
}

func (c RegexpConstraint) ValidateString(
	ctx context.Context,
	validator *validation.Validator,
	value *string,
) error {
	if c.regex == nil {
		return validator.CreateConstraintError("RegexpConstraint", "nil regex")
	}

	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}

	if c.match == c.regex.MatchString(*value) {
		return nil
	}

	return validator.
		BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		Create()
}

func IsJSON() validation.StringFuncConstraint {
	return validation.
		OfStringBy(predicate.JSON).
		WithError(validation.ErrInvalidJSON).
		WithMessage(validation.ErrInvalidJSON.Message())
}

func IsInteger() validation.StringFuncConstraint {
	return validation.
		OfStringBy(predicate.Integer).
		WithError(validation.ErrNotInteger).
		WithMessage(validation.ErrNotInteger.Message())
}

func IsNumeric() validation.StringFuncConstraint {
	return validation.
		OfStringBy(predicate.Number).
		WithError(validation.ErrNotNumeric).
		WithMessage(validation.ErrNotNumeric.Message())
}
