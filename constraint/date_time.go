package constraint

import (
	"context"
	"time"

	"line/validation"
)

type DateTimeConstraint struct {
	err               error
	layout            string
	messageTemplate   string
	groups            []string
	messageParameters validation.TemplateParameterList
	isIgnored         bool
}

func IsDateTime() DateTimeConstraint {
	return DateTimeConstraint{
		layout:          time.RFC3339,
		err:             validation.ErrInvalidDateTime,
		messageTemplate: validation.ErrInvalidDateTime.Message(),
	}
}

func IsDate() DateTimeConstraint {
	return DateTimeConstraint{
		layout:          "2006-01-02",
		err:             validation.ErrInvalidDate,
		messageTemplate: validation.ErrInvalidDate.Message(),
	}
}

func IsTime() DateTimeConstraint {
	return DateTimeConstraint{
		layout:          "15:04:05",
		err:             validation.ErrInvalidTime,
		messageTemplate: validation.ErrInvalidTime.Message(),
	}
}

func (c DateTimeConstraint) WithLayout(layout string) DateTimeConstraint {
	c.layout = layout
	return c
}

func (c DateTimeConstraint) WithError(err error) DateTimeConstraint {
	c.err = err
	return c
}

func (c DateTimeConstraint) WithMessage(
	template string,
	parameters ...validation.TemplateParameter,
) DateTimeConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters

	return c
}

func (c DateTimeConstraint) When(condition bool) DateTimeConstraint {
	c.isIgnored = !condition
	return c
}

func (c DateTimeConstraint) WhenGroups(groups ...string) DateTimeConstraint {
	c.groups = groups
	return c
}

func (c DateTimeConstraint) ValidateString(
	ctx context.Context,
	validator *validation.Validator,
	value *string,
) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}

	if _, err := time.Parse(c.layout, *value); err == nil {
		return nil
	}

	return validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ layout }}", Value: c.layout},
				validation.TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		WithParameter("{{ value }}", *value).Create()
}
