package constraint

import (
	"context"
	"fmt"
	"strings"

	"line/validation"
)

type ChoiceConstraint[T comparable] struct {
	blank             T
	choices           map[T]bool
	choicesValue      string
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
	disallowBlank     bool
	isIgnored         bool
}

func IsOneOf[T comparable](values ...T) ChoiceConstraint[T] {
	choices := make(map[T]bool, len(values))
	for _, value := range values {
		choices[value] = true
	}

	s := strings.Builder{}

	for i, value := range values {
		if i > 0 {
			s.WriteString(", ")
		}

		s.WriteString(fmt.Sprint(value))
	}

	return ChoiceConstraint[T]{
		choices:         choices,
		choicesValue:    s.String(),
		err:             validation.ErrNoSuchChoice,
		messageTemplate: validation.ErrNoSuchChoice.Message(),
	}
}

func (c ChoiceConstraint[T]) WithoutBlank() ChoiceConstraint[T] {
	c.disallowBlank = true
	return c
}

func (c ChoiceConstraint[T]) WithError(err error) ChoiceConstraint[T] {
	c.err = err
	return c
}

func (c ChoiceConstraint[T]) WithMessage(
	template string,
	parameters ...validation.TemplateParameter,
) ChoiceConstraint[T] {
	c.messageTemplate = template
	c.messageParameters = parameters

	return c
}

func (c ChoiceConstraint[T]) When(condition bool) ChoiceConstraint[T] {
	c.isIgnored = !condition
	return c
}

func (c ChoiceConstraint[T]) WhenGroups(groups ...string) ChoiceConstraint[T] {
	c.groups = groups
	return c
}

func (c ChoiceConstraint[T]) ValidateNumber(
	ctx context.Context,
	validator *validation.Validator,
	value *T,
) error {
	return c.ValidateComparable(ctx, validator, value)
}

func (c ChoiceConstraint[T]) ValidateString(
	ctx context.Context,
	validator *validation.Validator,
	value *T,
) error {
	return c.ValidateComparable(ctx, validator, value)
}

func (c ChoiceConstraint[T]) ValidateComparable(
	ctx context.Context,
	validator *validation.Validator,
	value *T,
) error {
	if len(c.choices) == 0 {
		return validator.CreateConstraintError("ChoiceConstraint", "empty list of choices")
	}

	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil ||
		!c.disallowBlank && *value == c.blank {
		return nil
	}

	if c.choices[*value] {
		return nil
	}

	return validator.
		BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: fmt.Sprint(*value)},
				validation.TemplateParameter{Key: "{{ choices }}", Value: c.choicesValue},
			)...,
		).
		Create()
}
