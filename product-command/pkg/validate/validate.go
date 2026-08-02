package validate_errors

import "github.com/go-playground/validator/v10"

func BrandsValidateError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			f := e.Field()

			switch f {
			case "Name":
				errors[f] = "name must be between 3 and 50 characters."
			}
		}
	}

	return errors
}
