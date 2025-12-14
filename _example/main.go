package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"line/constraint"
	"line/validation"
)

type User struct {
	Email    string
	Age      int
	Tags     []string
	IsActive bool
	Birthday string
}

func (u User) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx,

		validation.StringProperty(
			"email",
			u.Email,
			constraint.IsNotBlank(),
			constraint.HasMaxLength(255),
			constraint.Matches(regexp.MustCompile(`.+@.+\..+`)),
		),

		validation.ComparableProperty(
			"age",
			u.Age,
			constraint.IsNotBlankComparable[int](),
			constraint.IsOneOf(18, 21, 30, 40),
		),

		validation.EachStringProperty(
			"tags",
			u.Tags,
			constraint.HasMinLength(2),
			constraint.HasMaxLength(10),
		),

		validation.BoolProperty(
			"isActive",
			u.IsActive,
			constraint.IsNotBlankComparable[bool](),
		),

		validation.StringProperty(
			"birthday",
			u.Birthday,
			constraint.IsDate(),
		),
	)
}

func main() {
	validator, err := validation.NewValidator()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	user := User{
		Email:    "invalid-email",
		Age:      17,
		Tags:     []string{"a", "veryverylongtag"},
		IsActive: false,
		Birthday: "2023-99-99",
	}

	err = validator.Validate(ctx,

		validation.String(
			"he 021ddll",
			constraint.IsNotBlank(),
			constraint.Matches(regexp.MustCompile(`^he\s+\w+$`)),
		),

		validation.NilString(
			nil,
			constraint.IsNotNil(),
		),

		validation.Comparable(
			10,
			constraint.IsNotBlankNumber[int](),
		),

		validation.Countable(
			3,
			constraint.HasMinCount(5),
		),

		validation.Time(
			time.Now(),
			constraint.IsNotBlankComparable[time.Time](),
		),

		validation.String(
			`{invalid}`,
			constraint.IsJSON(),
		),

		validation.String(
			"123.45",
			constraint.IsNumeric(),
		),

		validation.String(
			"abc",
			constraint.IsInteger(),
		),

		validation.String(
			"blue",
			constraint.IsOneOf("red", "green"),
		),

		validation.When(user.IsActive).
			Then(
				validation.StringProperty(
					"email",
					user.Email,
					constraint.Matches(regexp.MustCompile(`@company\.com$`)),
				),
			).
			Else(
				validation.StringProperty(
					"email",
					user.Email,
					constraint.HasMaxLength(100),
				),
			),

		validation.String(
			"admin",
			constraint.
				IsOneOf("admin", "user").
				WhenGroups("admin"),
		),

		validation.All(
			validation.String(
				"",
				constraint.IsNotBlank(),
			),
			validation.String(
				"abc",
				constraint.HasMinLength(5),
			),
		),

		validation.AtLeastOneOf(
			validation.String(
				"",
				constraint.IsNotBlank(),
			),
			validation.String(
				"ok",
				constraint.HasMinLength(2),
			),
		),

		validation.Sequentially(
			validation.String(
				"",
				constraint.IsNotBlank(),
			),
			validation.String(
				"abc",
				constraint.HasMinLength(5),
			),
		),

		validation.Async(
			validation.String(
				"a",
				constraint.HasMinLength(2),
			),
			validation.String(
				"b",
				constraint.HasMinLength(2),
			),
		),

		validation.Valid(user),
	)

	if err != nil {
		if violations, ok := validation.UnwrapViolationList(err); ok {
			for _, v := range violations.AsSlice() {
				fmt.Println(v)
			}
		} else {
			log.Println(err)
		}
	}

}
