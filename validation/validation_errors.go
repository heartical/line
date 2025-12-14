package validation

import (
	"fmt"
	"strings"

	"line/message"
)

var (
	ErrInvalidDate       = NewError("invalid date", message.InvalidDate)
	ErrInvalidDateTime   = NewError("invalid datetime", message.InvalidDateTime)
	ErrInvalidJSON       = NewError("invalid JSON", message.InvalidJSON)
	ErrInvalidTime       = NewError("invalid time", message.InvalidTime)
	ErrIsBlank           = NewError("is blank", message.IsBlank)
	ErrIsEqual           = NewError("is equal", message.IsEqual)
	ErrIsNil             = NewError("is nil", message.IsNil)
	ErrNoSuchChoice      = NewError("no such choice", message.NoSuchChoice)
	ErrNotBlank          = NewError("is not blank", message.NotBlank)
	ErrNotDivisible      = NewError("is not divisible", message.NotDivisible)
	ErrNotDivisibleCount = NewError("not divisible count", message.NotDivisibleCount)
	ErrNotEqual          = NewError("is not equal", message.NotEqual)
	ErrNotExactCount     = NewError("not exact count", message.NotExactCount)
	ErrNotExactLength    = NewError("not exact length", message.NotExactLength)
	ErrNotFalse          = NewError("is not false", message.NotFalse)
	ErrNotInRange        = NewError("is not in range", message.NotInRange)
	ErrNotInteger        = NewError("is not an integer", message.NotInteger)
	ErrNotNegative       = NewError("is not negative", message.NotNegative)
	ErrNotNegativeOrZero = NewError("is not negative or zero", message.NotNegativeOrZero)
	ErrNotNil            = NewError("is not nil", message.NotNil)
	ErrNotNumeric        = NewError("is not numeric", message.NotNumeric)
	ErrNotPositive       = NewError("is not positive", message.NotPositive)
	ErrNotPositiveOrZero = NewError("is not positive or zero", message.NotPositiveOrZero)
	ErrNotTrue           = NewError("is not true", message.NotTrue)
	ErrNotUnique         = NewError("is not unique", message.NotUnique)
	ErrNotValid          = NewError("is not valid", message.NotValid)
	ErrProhibitedIP      = NewError("is prohibited IP", message.ProhibitedIP)
	ErrProhibitedURL     = NewError("is prohibited URL", message.ProhibitedURL)
	ErrTooEarly          = NewError("is too early", message.TooEarly)
	ErrTooEarlyOrEqual   = NewError("is too early or equal", message.TooEarlyOrEqual)
	ErrTooFewElements    = NewError("too few elements", message.TooFewElements)
	ErrTooHigh           = NewError("is too high", message.TooHigh)
	ErrTooHighOrEqual    = NewError("is too high or equal", message.TooHighOrEqual)
	ErrTooLate           = NewError("is too late", message.TooLate)
	ErrTooLateOrEqual    = NewError("is too late or equal", message.TooLateOrEqual)
	ErrTooLong           = NewError("is too long", message.TooLong)
	ErrTooLow            = NewError("is too low", message.TooLow)
	ErrTooLowOrEqual     = NewError("is too low or equal", message.TooLowOrEqual)
	ErrTooManyElements   = NewError("too many elements", message.TooManyElements)
	ErrTooShort          = NewError("is too short", message.TooShort)
)

type Error struct {
	code    string
	message string
}

func NewError(code, message string) *Error {
	return &Error{code: code, message: message}
}

func (err *Error) Error() string { return err.code }

func (err *Error) Message() string { return err.message }

type ConstraintError struct {
	ConstraintName string
	Path           *PropertyPath
	Description    string
}

func (err *ConstraintError) Error() string {
	var s strings.Builder

	s.WriteString("validate by " + err.ConstraintName)

	if err.Path != nil {
		s.WriteString(` at path "` + err.Path.String() + `"`)
	}

	s.WriteString(": " + err.Description)

	return s.String()
}

type ConstraintNotFoundError struct {
	Key  string
	Type string
}

func (err *ConstraintNotFoundError) Error() string {
	return fmt.Sprintf(
		"constraint by key %q of type %q is not found",
		err.Key,
		err.Type,
	)
}
