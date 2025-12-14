package validation

import (
	"context"
	"sync"
)

type WhenArgument struct {
	path          []PropertyPathElement
	thenArguments []Argument
	elseArguments []Argument
	isTrue        bool
}

func When(isTrue bool) WhenArgument {
	return WhenArgument{isTrue: isTrue}
}

func (arg WhenArgument) Then(arguments ...Argument) WhenArgument {
	arg.thenArguments = arguments
	return arg
}

func (arg WhenArgument) Else(arguments ...Argument) WhenArgument {
	arg.elseArguments = arguments
	return arg
}

func (arg WhenArgument) At(path ...PropertyPathElement) WhenArgument {
	arg.path = append(arg.path, path...)
	return arg
}

func (arg WhenArgument) setUp(ctx *executionContext) {
	ctx.addValidation(arg.validate, arg.path...)
}

func (arg WhenArgument) validate(
	ctx context.Context,
	validator *Validator,
) (*ViolationListError, error) {
	var err error
	if arg.isTrue {
		err = validator.Validate(ctx, arg.thenArguments...)
	} else {
		err = validator.Validate(ctx, arg.elseArguments...)
	}

	return unwrapViolationList(err)
}

type WhenGroupsArgument struct {
	groups        []string
	path          []PropertyPathElement
	thenArguments []Argument
	elseArguments []Argument
}

func WhenGroups(groups ...string) WhenGroupsArgument {
	return WhenGroupsArgument{groups: groups}
}

func (arg WhenGroupsArgument) Then(arguments ...Argument) WhenGroupsArgument {
	arg.thenArguments = arguments
	return arg
}

func (arg WhenGroupsArgument) Else(arguments ...Argument) WhenGroupsArgument {
	arg.elseArguments = arguments
	return arg
}

func (arg WhenGroupsArgument) At(path ...PropertyPathElement) WhenGroupsArgument {
	arg.path = append(arg.path, path...)
	return arg
}

func (arg WhenGroupsArgument) setUp(ctx *executionContext) {
	ctx.addValidation(arg.validate, arg.path...)
}

func (arg WhenGroupsArgument) validate(
	ctx context.Context,
	validator *Validator,
) (*ViolationListError, error) {
	var err error
	if validator.IsIgnoredForGroups(arg.groups...) {
		err = validator.Validate(ctx, arg.elseArguments...)
	} else {
		err = validator.Validate(ctx, arg.thenArguments...)
	}

	return unwrapViolationList(err)
}

type SequentialArgument struct {
	path      []PropertyPathElement
	arguments []Argument
	isIgnored bool
}

func Sequentially(arguments ...Argument) SequentialArgument {
	return SequentialArgument{arguments: arguments}
}

func (arg SequentialArgument) At(path ...PropertyPathElement) SequentialArgument {
	arg.path = append(arg.path, path...)
	return arg
}

func (arg SequentialArgument) When(condition bool) SequentialArgument {
	arg.isIgnored = !condition
	return arg
}

func (arg SequentialArgument) setUp(ctx *executionContext) {
	ctx.addValidation(arg.validate, arg.path...)
}

func (arg SequentialArgument) validate(
	ctx context.Context,
	validator *Validator,
) (*ViolationListError, error) {
	if arg.isIgnored {
		return &ViolationListError{}, nil
	}

	violations := &ViolationListError{}

	for _, argument := range arg.arguments {
		err := violations.AppendFromError(validator.Validate(ctx, argument))
		if err != nil {
			return nil, err
		}

		if violations.len > 0 {
			return violations, nil
		}
	}

	return violations, nil
}

type AtLeastOneOfArgument struct {
	path      []PropertyPathElement
	arguments []Argument
	isIgnored bool
}

func AtLeastOneOf(arguments ...Argument) AtLeastOneOfArgument {
	return AtLeastOneOfArgument{arguments: arguments}
}

func (arg AtLeastOneOfArgument) At(path ...PropertyPathElement) AtLeastOneOfArgument {
	arg.path = append(arg.path, path...)
	return arg
}

func (arg AtLeastOneOfArgument) When(condition bool) AtLeastOneOfArgument {
	arg.isIgnored = !condition
	return arg
}

func (arg AtLeastOneOfArgument) setUp(ctx *executionContext) {
	ctx.addValidation(arg.validate, arg.path...)
}

func (arg AtLeastOneOfArgument) validate(
	ctx context.Context,
	validator *Validator,
) (*ViolationListError, error) {
	if arg.isIgnored {
		return &ViolationListError{}, nil
	}

	violations := &ViolationListError{}

	for _, argument := range arg.arguments {
		violation := validator.Validate(ctx, argument)
		if violation == nil {
			return &ViolationListError{}, nil
		}

		err := violations.AppendFromError(violation)
		if err != nil {
			return nil, err
		}
	}

	return violations, nil
}

type AllArgument struct {
	path      []PropertyPathElement
	arguments []Argument
	isIgnored bool
}

func All(arguments ...Argument) AllArgument {
	return AllArgument{arguments: arguments}
}

func AtProperty(propertyName string, arguments ...Argument) AllArgument {
	return All(arguments...).At(PropertyName(propertyName))
}

func (arg AllArgument) At(path ...PropertyPathElement) AllArgument {
	arg.path = append(arg.path, path...)
	return arg
}

func (arg AllArgument) When(condition bool) AllArgument {
	arg.isIgnored = !condition
	return arg
}

func (arg AllArgument) setUp(ctx *executionContext) {
	ctx.addValidation(arg.validate, arg.path...)
}

func (arg AllArgument) validate(
	ctx context.Context,
	validator *Validator,
) (*ViolationListError, error) {
	if arg.isIgnored {
		return &ViolationListError{}, nil
	}

	violations := &ViolationListError{}

	for _, argument := range arg.arguments {
		err := violations.AppendFromError(validator.Validate(ctx, argument))
		if err != nil {
			return nil, err
		}
	}

	return violations, nil
}

type AsyncArgument struct {
	path      []PropertyPathElement
	arguments []Argument
	isIgnored bool
}

func Async(arguments ...Argument) AsyncArgument {
	return AsyncArgument{arguments: arguments}
}

func (arg AsyncArgument) At(path ...PropertyPathElement) AsyncArgument {
	arg.path = append(arg.path, path...)
	return arg
}

func (arg AsyncArgument) When(condition bool) AsyncArgument {
	arg.isIgnored = !condition
	return arg
}

func (arg AsyncArgument) setUp(ctx *executionContext) {
	ctx.addValidation(arg.validate, arg.path...)
}

func (arg AsyncArgument) validate(
	ctx context.Context,
	validator *Validator,
) (*ViolationListError, error) {
	if arg.isIgnored {
		return &ViolationListError{}, nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	waiter := &sync.WaitGroup{}
	waiter.Add(len(arg.arguments))

	errs := make(chan error)

	for _, argument := range arg.arguments {
		go func(argument Argument) {
			defer waiter.Done()

			errs <- validator.Validate(ctx, argument)
		}(argument)
	}

	go func() {
		waiter.Wait()
		close(errs)
	}()

	violations := &ViolationListError{}

	for violation := range errs {
		err := violations.AppendFromError(violation)
		if err != nil {
			return nil, err
		}
	}

	return violations, nil
}
