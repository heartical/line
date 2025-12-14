package validation

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type PropertyPathElement interface {
	IsIndex() bool
	fmt.Stringer
}

type PropertyName string

func (p PropertyName) IsIndex() bool {
	return false
}

func (p PropertyName) String() string {
	return string(p)
}

type ArrayIndex int

func (a ArrayIndex) IsIndex() bool {
	return true
}

func (a ArrayIndex) String() string {
	return strconv.Itoa(int(a))
}

type PropertyPath struct {
	parent *PropertyPath
	value  PropertyPathElement
}

func NewPropertyPath(elements ...PropertyPathElement) *PropertyPath {
	var path *PropertyPath

	return path.With(elements...)
}

func (path *PropertyPath) With(elements ...PropertyPathElement) *PropertyPath {
	current := path
	for _, element := range elements {
		current = &PropertyPath{parent: current, value: element}
	}

	return current
}

func (path *PropertyPath) WithProperty(name string) *PropertyPath {
	return &PropertyPath{
		parent: path,
		value:  PropertyName(name),
	}
}

func (path *PropertyPath) WithIndex(index int) *PropertyPath {
	return &PropertyPath{
		parent: path,
		value:  ArrayIndex(index),
	}
}

func (path *PropertyPath) Elements() []PropertyPathElement {
	if path == nil || path.value == nil {
		return nil
	}

	length := path.Len()
	elements := make([]PropertyPathElement, length)

	i := length - 1

	element := path
	for element != nil {
		elements[i] = element.value
		element = element.parent
		i--
	}

	return elements
}

func (path *PropertyPath) Len() int {
	length := 0
	element := path

	for element != nil {
		length++
		element = element.parent
	}

	return length
}

func (path *PropertyPath) String() string {
	elements := path.Elements()
	count := 0

	for _, element := range elements {
		if s, ok := element.(PropertyName); ok {
			count += len(s)
		} else {
			count += 2
		}
	}

	s := strings.Builder{}
	s.Grow(count)

	for i, element := range elements {
		name := element.String()

		switch {
		case element.IsIndex():
			s.WriteString("[" + name + "]")
		case isIdentifier(name):
			if i > 0 {
				s.WriteString(".")
			}

			s.WriteString(name)
		default:
			s.WriteString("['")
			writePropertyName(&s, name)
			s.WriteString("']")
		}
	}

	return s.String()
}

func (path *PropertyPath) MarshalText() ([]byte, error) {
	return []byte(path.String()), nil
}

func (path *PropertyPath) UnmarshalText(text []byte) error {
	parser := pathParser{}

	p, err := parser.Parse(string(text))
	if p == nil || err != nil {
		return err
	}

	*path = *p

	return nil
}

func isIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}

	for i, c := range s {
		if i == 0 && !isFirstIdentifierChar(c) {
			return false
		}

		if i > 0 && !isIdentifierChar(c) {
			return false
		}
	}

	return true
}

func isFirstIdentifierChar(c rune) bool {
	return unicode.IsLetter(c) || c == '$' || c == '_'
}

func isIdentifierChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '$' || c == '_'
}

func writePropertyName(s *strings.Builder, name string) {
	for _, c := range name {
		if c == '\'' || c == '\\' {
			s.WriteByte('\\')
		}

		s.WriteRune(c)
	}
}

type parsingState byte

const (
	initialState parsingState = iota
	beginIdentifierState
	identifierState

	beginIndexState
	indexState

	bracketedNameState
	endBracketedNameState

	closeBracketState
)

type pathParser struct {
	path      *PropertyPath
	buffer    strings.Builder
	index     int
	pathIndex int
	state     parsingState
	isEscape  bool
}

func (parser *pathParser) Parse(encodedPath string) (*PropertyPath, error) {
	if len(encodedPath) == 0 {
		return &PropertyPath{}, nil
	}

	for i, c := range encodedPath {
		parser.index = i
		if err := parser.handleNext(c); err != nil {
			return nil, err
		}
	}

	return parser.finish()
}

func (parser *pathParser) handleNext(c rune) error {
	var err error

	switch c {
	case '[':
		err = parser.handleOpenBracket(c)
	case ']':
		err = parser.handleCloseBracket(c)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		err = parser.handleDigit(c)
	case '\'':
		err = parser.handleQuote(c)
	case '\\':
		err = parser.handleEscape(c)
	case '.':
		err = parser.handlePoint(c)
	default:
		err = parser.handleOther(c)
	}

	return err
}

func (parser *pathParser) handleOpenBracket(c rune) error {
	switch parser.state {
	case beginIdentifierState, beginIndexState, indexState, endBracketedNameState:
		return parser.newCharError(c, "unexpected char")
	case identifierState:
		if parser.buffer.Len() > 0 {
			parser.addProperty()
		}

		parser.state = beginIndexState
	case initialState, closeBracketState:
		parser.state = beginIndexState
	case bracketedNameState:
		parser.buffer.WriteRune(c)
	}

	return nil
}

func (parser *pathParser) handleCloseBracket(c rune) error {
	switch parser.state {
	case indexState:
		err := parser.addIndex()
		if err != nil {
			return err
		}

		parser.state = closeBracketState
	case bracketedNameState:
		parser.buffer.WriteRune(c)
	case endBracketedNameState:
		parser.addProperty()
		parser.state = closeBracketState
	default:
		return parser.newCharError(c, "unexpected close bracket")
	}

	return nil
}

func (parser *pathParser) handleDigit(c rune) error {
	switch parser.state {
	case beginIndexState, indexState:
		parser.state = indexState
	case bracketedNameState, identifierState:
	case initialState, beginIdentifierState:
		return parser.newCharError(c, "unexpected identifier character")
	default:
		return parser.newCharError(c, "invalid array index")
	}

	parser.buffer.WriteRune(c)

	return nil
}

func (parser *pathParser) handlePoint(c rune) error {
	switch parser.state {
	case beginIdentifierState, identifierState:
		if parser.buffer.Len() == 0 {
			return parser.newCharError(c, "unexpected point")
		}

		parser.addProperty()
		parser.state = beginIdentifierState
	case bracketedNameState:
		parser.buffer.WriteRune(c)
	case closeBracketState:
		parser.state = beginIdentifierState
	default:
		return parser.newCharError(c, "unexpected point")
	}

	return nil
}

func (parser *pathParser) handleQuote(c rune) error {
	if parser.isEscape {
		parser.buffer.WriteRune(c)
		parser.isEscape = false

		return nil
	}

	switch parser.state {
	case beginIndexState:
		parser.state = bracketedNameState
	case bracketedNameState:
		parser.state = endBracketedNameState
	default:
		return parser.newCharError(c, "unexpected quote")
	}

	return nil
}

func (parser *pathParser) handleEscape(c rune) error {
	if parser.state != bracketedNameState {
		return parser.newCharError(c, "unexpected backslash")
	}

	if parser.isEscape {
		parser.buffer.WriteRune(c)
		parser.isEscape = false
	} else {
		parser.isEscape = true
	}

	return nil
}

func (parser *pathParser) handleOther(c rune) error {
	switch parser.state {
	case beginIndexState, indexState:
		return parser.newCharError(c, "unexpected array index character")
	case initialState, beginIdentifierState, identifierState:
		if !isFirstIdentifierChar(c) {
			return parser.newCharError(c, "unexpected identifier char")
		}

		parser.state = identifierState
	case closeBracketState, endBracketedNameState:
		return parser.newCharError(c, "unexpected char")
	}

	parser.buffer.WriteRune(c)

	return nil
}

func (parser *pathParser) addProperty() {
	parser.path = parser.path.WithProperty(parser.buffer.String())
	parser.pathIndex++
	parser.buffer.Reset()
}

func (parser *pathParser) addIndex() error {
	s := parser.buffer.String()

	u, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		if errors.Is(err, strconv.ErrRange) {
			return parser.newProcessingError("value out of range: " + s)
		}

		return parser.newProcessingError("invalid array index: " + s)
	}

	if u > math.MaxInt {
		return parser.newProcessingError("value out of range: " + s)
	}

	parser.path = parser.path.WithIndex(int(u))
	parser.pathIndex++
	parser.buffer.Reset()

	return nil
}

func (parser *pathParser) finish() (*PropertyPath, error) {
	switch parser.state {
	case beginIdentifierState, identifierState:
		if parser.buffer.Len() == 0 {
			return nil, parser.newError("incomplete property name")
		}

		parser.path = parser.path.WithProperty(parser.buffer.String())
	case beginIndexState, indexState:
		return nil, parser.newError("incomplete array index")
	case bracketedNameState, endBracketedNameState:
		return nil, parser.newError("incomplete bracketed property name")
	case closeBracketState:
	default:
		return nil, parser.newError("unexpected parsing state")
	}

	return parser.path, nil
}

func (parser *pathParser) newError(message string) *pathParsingError {
	return &pathParsingError{
		pathIndex: parser.pathIndex,
		message:   message,
	}
}

func (parser *pathParser) newCharError(char rune, message string) *pathParsingCharError {
	return &pathParsingCharError{
		index:     parser.index,
		pathIndex: parser.pathIndex,
		char:      char,
		message:   message,
	}
}

func (parser *pathParser) newProcessingError(message string) *pathParsingProcessingError {
	return &pathParsingProcessingError{
		pathIndex: parser.pathIndex,
		message:   message,
	}
}

type pathParsingError struct {
	message   string
	pathIndex int
}

func (err *pathParsingError) Error() string {
	return fmt.Sprintf("parsing path element #%d: %s", err.pathIndex, err.message)
}

type pathParsingCharError struct {
	message   string
	index     int
	pathIndex int
	char      rune
}

func (err *pathParsingCharError) Error() string {
	return fmt.Sprintf(
		"parsing path element #%d at char #%d %q: %s",
		err.pathIndex,
		err.index,
		err.char,
		err.message,
	)
}

type pathParsingProcessingError struct {
	message   string
	pathIndex int
}

func (err *pathParsingProcessingError) Error() string {
	return fmt.Sprintf("parsing path element #%d: %s", err.pathIndex, err.message)
}
