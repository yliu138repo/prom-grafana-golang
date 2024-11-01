package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var DtoValidator *validator.Validate

var uuidV4Regex = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

func uuidv4(fl validator.FieldLevel) bool {
	return uuidV4Regex.MatchString(fl.Field().String())
}

func init() {
	DtoValidator = validator.New()
	DtoValidator.RegisterValidation("uuid4", uuidv4)
}

func ValidateStruct(s interface{}) []ErrorResponse {
	var errors []ErrorResponse

	err := DtoValidator.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var msg string

			switch err.Tag() {
			case "required":
				msg = "This field is required."
			case "min":
				msg = "Value is too short."
			case "max":
				msg = "Value is too long."
			case "email":
				msg = "Invalid email format."
			case "gte":
				msg = "Value must be greater than or equal to required threshold."
			case "lte":
				msg = "Value must be less than or equal to required threshold."
			// Add more cases as needed
			default:
				msg = "Invalid value."
			}

			errors = append(errors, ErrorResponse{
				Field:   err.Field(),
				Message: msg,
			})
		}
	}

	return errors
}
