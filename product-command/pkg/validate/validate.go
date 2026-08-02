package validate

import "github.com/go-playground/validator/v10"

func BrandsValidateError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			f := e.Field()
			tag := e.Tag()

			switch f {
			case "Name":
				switch tag {
				case "required":
					errors[f] = "name is required"
				}
			}
		}
	}
}
