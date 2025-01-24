package auth

import (
	"errors"
	"regexp"
	"sync"
	"atorgold/response"
	"github.com/go-playground/validator/v10"
)

var customValidators = map[string]func(validator.FieldLevel) bool{
	"email": func(fl validator.FieldLevel) bool {
		// Basic email regex
		email := fl.Field().String()
		regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
		return regex.MatchString(email)
	},
	"password": func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		//var isValid bool
		isValid := true
		switch {
		case regexp.MustCompile(".*[A-Z].*").MatchString(password):
			fallthrough
		case regexp.MustCompile(".*[a-z].*").MatchString(password):
			fallthrough
		case regexp.MustCompile(".*\\d.*").MatchString(password):
			fallthrough
		case regexp.MustCompile(".*[@*#$%^&+=!].*").MatchString(password):
			fallthrough
		case regexp.MustCompile(".{8,20}").MatchString(password):
			isValid = true
		}
		return isValid
	},
	"name": func(fl validator.FieldLevel) bool {
		name := fl.Field().String()
		regex := regexp.MustCompile("^[a-zA-Z ]+$")
		return regex.MatchString(name)
	},
	"surname": func(fl validator.FieldLevel) bool {
		surname := fl.Field().String()
		regex := regexp.MustCompile("^[a-zA-Z ]+$")
		return regex.MatchString(surname)
	},
	"phone": func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		regex := regexp.MustCompile("^[0-9]{10}$")
		return regex.MatchString(phone)
	},
	"image_buffer_1": func(fl validator.FieldLevel) bool {
		imageBuffer := fl.Field().String()
		return len(imageBuffer) > 0
	},
	"image_name_1": func(fl validator.FieldLevel) bool {
		imageName := fl.Field().String()
		return len(imageName) > 0
	},
}

var validations *validator.Validate
var once sync.Once

func CreateNewValidator() *validator.Validate {
	once.Do(func() {
		validations = validator.New()
		for key, value := range customValidators {
			validations.RegisterValidation(key, value)
		}
	})
	return validations
}

func ValidateRequestByStruct[T any](s T) []*response.ValidationMessage {
	validate := CreateNewValidator()
	var allErrors []*response.ValidationMessage
	err := validate.Struct(s)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			allErrors = append(allErrors, &response.ValidationMessage{
				FailedField: "N/A",
				Tag:         "invalid",
				Message:     err.Error(),
			})
			return allErrors
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, err := range validationErrors {
				var element response.ValidationMessage
				element.FailedField = err.Field()
				element.Tag = err.Tag()
				element.Message = err.Error()
				allErrors = append(allErrors, &element)
			}
		}
	}
	return allErrors
}
