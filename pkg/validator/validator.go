package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
	// Register custom validation for gender
	validate.RegisterValidation("gender", validateGender)
}

func Validate(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		return err
	}
	return nil
}

func validateGender(fl validator.FieldLevel) bool {
	gender := strings.ToLower(fl.Field().String())
	return gender == "male" || gender == "female"
}

// ErrorMap converts validation errors into a map
func ErrorMap(err error) map[string]string {
	if err == nil {
		return nil
	}

	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return map[string]string{"general": "validation error"}
	}

	fieldErrors := make(map[string]string)
	for _, e := range errs {
		field := e.Field()
		msg := "invalid value"

		switch e.Tag() {
		case "required":
			msg = "field is required"
		case "email":
			msg = "must be a valid email"
		case "min":
			msg = fmt.Sprintf("must be at least %v characters", e.Param())
		case "max":
			msg = fmt.Sprintf("must be at most %v characters", e.Param())
		case "gender":
			msg = "must be male or female"
		}

		fieldErrors[field] = msg
	}

	return fieldErrors
}
