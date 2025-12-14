package validation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	initialBufferSize = 32
)

type Violation interface {
	error
	Unwrap() error
	Is(target error) bool
	Message() string
	MessageTemplate() string
	Parameters() []TemplateParameter
	PropertyPath() *PropertyPath
}

type ViolationFactory interface {
	CreateViolation(
		err error,
		messageTemplate string,
		parameters []TemplateParameter,
		propertyPath *PropertyPath,
	) Violation
}

type NewViolationFunc func(
	err error,
	messageTemplate string,
	parameters []TemplateParameter,
	propertyPath *PropertyPath,
) Violation

func (f NewViolationFunc) CreateViolation(
	err error,
	messageTemplate string,
	parameters []TemplateParameter,
	propertyPath *PropertyPath,
) Violation {
	return f(err, messageTemplate, parameters, propertyPath)
}

type ViolationListError struct {
	first *ViolationListElementError
	last  *ViolationListElementError
	len   int
}

type ViolationListElementError struct {
	next      *ViolationListElementError
	violation Violation
}

func NewViolationList(violations ...Violation) *ViolationListError {
	list := &ViolationListError{}
	list.Append(violations...)

	return list
}

func (list *ViolationListError) Len() int {
	if list == nil {
		return 0
	}

	return list.len
}

func (list *ViolationListError) ForEach(f func(i int, violation Violation) error) error {
	if list == nil {
		return nil
	}

	i := 0
	for e := list.first; e != nil; e = e.next {
		err := f(i, e.violation)
		if err != nil {
			return err
		}

		i++
	}

	return nil
}

func (list *ViolationListError) First() *ViolationListElementError {
	return list.first
}

func (list *ViolationListError) Last() *ViolationListElementError {
	return list.last
}

func (list *ViolationListError) Append(violations ...Violation) {
	for i := range violations {
		element := &ViolationListElementError{violation: violations[i]}
		if list.first == nil {
			list.first = element
			list.last = element
		} else {
			list.last.next = element
			list.last = element
		}
	}

	list.len += len(violations)
}

func (list *ViolationListError) Join(violations *ViolationListError) {
	if violations == nil || violations.len == 0 {
		return
	}

	if list.first == nil {
		list.first = violations.first
		list.last = violations.last
	} else {
		list.last.next = violations.first
		list.last = violations.last
	}

	list.len += violations.len
}

func (list *ViolationListError) Error() string {
	if list == nil || list.len == 0 {
		return "the list of violations is empty, it looks like you forgot to use the AsError method somewhere"
	}

	return list.String()
}

func (list *ViolationListError) String() string {
	return list.toString(" ")
}

func (list *ViolationListError) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = io.WriteString(f, list.toString("\n\t"))
		} else {
			_, _ = io.WriteString(f, list.toString(" "))
		}
	case 's', 'q':
		_, _ = io.WriteString(f, list.toString(" "))
	}
}

func (list *ViolationListError) toString(delimiter string) string {
	if list == nil || list.len == 0 {
		return ""
	}

	if list.len == 1 {
		return list.first.violation.Error()
	}

	var s strings.Builder

	s.Grow(initialBufferSize * list.len)
	s.WriteString("violations:")

	i := 0

	for e := list.first; e != nil; e = e.next {
		v := e.violation

		if i > 0 {
			s.WriteString(";")
		}

		s.WriteString(delimiter)
		s.WriteString("#" + strconv.Itoa(i))

		if v.PropertyPath() != nil {
			s.WriteString(` at "` + v.PropertyPath().String() + `"`)
		}

		s.WriteString(`: "` + v.Message() + `"`)

		i++
	}

	return s.String()
}

func (list *ViolationListError) AppendFromError(err error) error {
	if err == nil {
		return nil
	}

	if violation, ok := UnwrapViolation(err); ok {
		list.Append(violation)
		return nil
	}

	if violationList, ok := UnwrapViolationList(err); ok {
		list.Join(violationList)
		return nil
	}

	return err
}

func (list *ViolationListError) Is(target error) bool {
	for e := list.first; e != nil; e = e.next {
		if e.violation.Is(target) {
			return true
		}
	}

	return false
}

func (list *ViolationListError) Filter(errs ...error) *ViolationListError {
	filtered := &ViolationListError{}

	for e := list.first; e != nil; e = e.next {
		for _, err := range errs {
			if e.violation.Is(err) {
				filtered.Append(e.violation)
			}
		}
	}

	return filtered
}

func (list *ViolationListError) AsError() error {
	if list == nil || list.len == 0 {
		return nil
	}

	return list
}

func (list *ViolationListError) AsSlice() []Violation {
	violations := make([]Violation, list.len)

	i := 0

	for e := list.first; e != nil; e = e.next {
		violations[i] = e.violation
		i++
	}

	return violations
}

func (list *ViolationListError) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	b.WriteRune('[')

	i := 0

	for e := list.first; e != nil; e = e.next {
		data, err := json.Marshal(e.violation)
		if err != nil {
			return nil, fmt.Errorf("marshal violation at %d: %w", i, err)
		}

		b.Write(data)

		if e.next != nil {
			b.WriteRune(',')
		}

		i++
	}

	b.WriteRune(']')

	return b.Bytes(), nil
}

func (element *ViolationListElementError) Next() *ViolationListElementError {
	return element.next
}

func (element *ViolationListElementError) Violation() Violation {
	return element.violation
}

func (element *ViolationListElementError) Unwrap() error {
	return element.violation.Unwrap()
}

func (element *ViolationListElementError) Error() string {
	return element.violation.Error()
}

func (element *ViolationListElementError) Is(target error) bool {
	return element.violation.Is(target)
}

func (element *ViolationListElementError) Message() string {
	return element.violation.Message()
}

func (element *ViolationListElementError) MessageTemplate() string {
	return element.violation.MessageTemplate()
}

func (element *ViolationListElementError) Parameters() []TemplateParameter {
	return element.violation.Parameters()
}

func (element *ViolationListElementError) PropertyPath() *PropertyPath {
	return element.violation.PropertyPath()
}

func IsViolation(err error) bool {
	var violation Violation

	return errors.As(err, &violation)
}

func IsViolationList(err error) bool {
	var violations *ViolationListError

	return errors.As(err, &violations)
}

func UnwrapViolation(err error) (Violation, bool) {
	var violation Violation

	as := errors.As(err, &violation)

	return violation, as
}

func UnwrapViolationList(err error) (*ViolationListError, bool) {
	var violations *ViolationListError

	as := errors.As(err, &violations)

	return violations, as
}

type internalViolationError struct {
	err             error
	propertyPath    *PropertyPath
	message         string
	messageTemplate string
	parameters      []TemplateParameter
}

func (v *internalViolationError) Unwrap() error {
	return v.err
}

func (v *internalViolationError) Is(target error) bool {
	return errors.Is(v.err, target)
}

func (v *internalViolationError) Error() string {
	var s strings.Builder

	s.Grow(initialBufferSize)
	v.writeToBuilder(&s)

	return s.String()
}

func (v *internalViolationError) writeToBuilder(s *strings.Builder) {
	s.WriteString("violation")

	if v.propertyPath != nil {
		s.WriteString(` at "` + v.propertyPath.String() + `"`)
	}

	s.WriteString(`: "` + v.message + `"`)
}

func (v *internalViolationError) Message() string { return v.message }

func (v *internalViolationError) MessageTemplate() string { return v.messageTemplate }

func (v *internalViolationError) Parameters() []TemplateParameter { return v.parameters }

func (v *internalViolationError) PropertyPath() *PropertyPath { return v.propertyPath }

func (v *internalViolationError) MarshalJSON() ([]byte, error) {
	data := struct {
		PropertyPath *PropertyPath `json:"propertyPath,omitempty"`
		Error        string        `json:"error,omitempty"`
		Message      string        `json:"message"`
	}{
		Message:      v.message,
		PropertyPath: v.propertyPath,
	}
	if v.err != nil {
		data.Error = v.err.Error()
	}

	return json.Marshal(data)
}

type BuiltinViolationFactory struct{}

func NewViolationFactory() *BuiltinViolationFactory {
	return &BuiltinViolationFactory{}
}

func (factory *BuiltinViolationFactory) CreateViolation(
	err error,
	messageTemplate string,
	parameters []TemplateParameter,
	propertyPath *PropertyPath,
) Violation {
	message := messageTemplate

	return &internalViolationError{
		err:             err,
		message:         renderMessage(message, parameters),
		messageTemplate: messageTemplate,
		parameters:      parameters,
		propertyPath:    propertyPath,
	}
}

type ViolationBuilder struct {
	err              error
	violationFactory ViolationFactory
	propertyPath     *PropertyPath
	messageTemplate  string
	parameters       []TemplateParameter
}

func NewViolationBuilder(factory ViolationFactory) *ViolationBuilder {
	return &ViolationBuilder{violationFactory: factory}
}

func (b *ViolationBuilder) BuildViolation(err error, message string) *ViolationBuilder {
	return &ViolationBuilder{
		err:              err,
		messageTemplate:  message,
		violationFactory: b.violationFactory,
	}
}

func (b *ViolationBuilder) SetPropertyPath(path *PropertyPath) *ViolationBuilder {
	b.propertyPath = path

	return b
}

func (b *ViolationBuilder) WithParameters(parameters ...TemplateParameter) *ViolationBuilder {
	b.parameters = parameters

	return b
}

func (b *ViolationBuilder) WithParameter(name, value string) *ViolationBuilder {
	b.parameters = append(b.parameters, TemplateParameter{Key: name, Value: value})

	return b
}

func (b *ViolationBuilder) At(path ...PropertyPathElement) *ViolationBuilder {
	b.propertyPath = b.propertyPath.With(path...)

	return b
}

func (b *ViolationBuilder) AtProperty(propertyName string) *ViolationBuilder {
	b.propertyPath = b.propertyPath.WithProperty(propertyName)

	return b
}

func (b *ViolationBuilder) AtIndex(index int) *ViolationBuilder {
	b.propertyPath = b.propertyPath.WithIndex(index)

	return b
}

func (b *ViolationBuilder) Create() Violation {
	return b.violationFactory.CreateViolation(
		b.err,
		b.messageTemplate,
		b.parameters,
		b.propertyPath,
	)
}

type ViolationListBuilder struct {
	violations       *ViolationListError
	violationFactory ViolationFactory

	propertyPath *PropertyPath
}

type ViolationListElementBuilder struct {
	err             error
	listBuilder     *ViolationListBuilder
	propertyPath    *PropertyPath
	messageTemplate string
	parameters      []TemplateParameter
}

func NewViolationListBuilder(factory ViolationFactory) *ViolationListBuilder {
	return &ViolationListBuilder{violationFactory: factory, violations: NewViolationList()}
}

func (b *ViolationListBuilder) BuildViolation(
	err error,
	message string,
) *ViolationListElementBuilder {
	return &ViolationListElementBuilder{
		listBuilder:     b,
		err:             err,
		messageTemplate: message,
		propertyPath:    b.propertyPath,
	}
}

func (b *ViolationListBuilder) AddViolation(
	err error,
	message string,
	path ...PropertyPathElement,
) *ViolationListBuilder {
	return b.add(err, message, nil, b.propertyPath.With(path...))
}

func (b *ViolationListBuilder) SetPropertyPath(path *PropertyPath) *ViolationListBuilder {
	b.propertyPath = path

	return b
}

func (b *ViolationListBuilder) At(path ...PropertyPathElement) *ViolationListBuilder {
	b.propertyPath = b.propertyPath.With(path...)

	return b
}

func (b *ViolationListBuilder) AtProperty(propertyName string) *ViolationListBuilder {
	b.propertyPath = b.propertyPath.WithProperty(propertyName)

	return b
}

func (b *ViolationListBuilder) AtIndex(index int) *ViolationListBuilder {
	b.propertyPath = b.propertyPath.WithIndex(index)

	return b
}

func (b *ViolationListBuilder) Create() *ViolationListError {
	return b.violations
}

func (b *ViolationListBuilder) add(
	err error,
	template string,
	parameters []TemplateParameter,
	path *PropertyPath,
) *ViolationListBuilder {
	b.violations.Append(b.violationFactory.CreateViolation(
		err,
		template,
		parameters,
		path,
	))

	return b
}

func (b *ViolationListElementBuilder) WithParameters(
	parameters ...TemplateParameter,
) *ViolationListElementBuilder {
	b.parameters = parameters

	return b
}

func (b *ViolationListElementBuilder) WithParameter(
	name, value string,
) *ViolationListElementBuilder {
	b.parameters = append(b.parameters, TemplateParameter{Key: name, Value: value})

	return b
}

func (b *ViolationListElementBuilder) At(path ...PropertyPathElement) *ViolationListElementBuilder {
	b.propertyPath = b.propertyPath.With(path...)

	return b
}

func (b *ViolationListElementBuilder) AtProperty(propertyName string) *ViolationListElementBuilder {
	b.propertyPath = b.propertyPath.WithProperty(propertyName)

	return b
}

func (b *ViolationListElementBuilder) AtIndex(index int) *ViolationListElementBuilder {
	b.propertyPath = b.propertyPath.WithIndex(index)

	return b
}

func (b *ViolationListElementBuilder) Add() *ViolationListBuilder {
	return b.listBuilder.add(b.err, b.messageTemplate, b.parameters, b.propertyPath)
}

func unwrapViolationList(err error) (*ViolationListError, error) {
	violations := NewViolationList()

	fatal := violations.AppendFromError(err)
	if fatal != nil {
		return nil, fatal
	}

	return violations, nil
}
