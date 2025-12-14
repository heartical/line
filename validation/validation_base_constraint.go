package validation

import "context"

type BaseConstraint struct {
	Err             error
	MessageTemplate string
	Groups          []string
	Parameters      TemplateParameterList
	IsIgnored       bool
}

func (c BaseConstraint) When(condition bool) BaseConstraint {
	c.IsIgnored = !condition
	return c
}

func (c BaseConstraint) WhenGroups(groups ...string) BaseConstraint {
	c.Groups = groups
	return c
}

func (c BaseConstraint) WithError(err error) BaseConstraint {
	c.Err = err
	return c
}

func (c BaseConstraint) WithMessage(
	template string,
	parameters ...TemplateParameter,
) BaseConstraint {
	c.MessageTemplate = template
	c.Parameters = parameters

	return c
}

func (c BaseConstraint) ShouldSkip(validator *Validator) bool {
	return c.IsIgnored || validator.IsIgnoredForGroups(c.Groups...)
}

func (c BaseConstraint) NewViolation(
	ctx context.Context,
	validator *Validator,
) Violation {
	return validator.
		BuildViolation(ctx, c.Err, c.MessageTemplate).
		WithParameters(c.Parameters...).
		Create()
}
