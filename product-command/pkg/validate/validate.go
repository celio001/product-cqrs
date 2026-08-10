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

func CategoriesValidateError(err error) map[string]string {
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

func ProductValidateError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			f := e.Field()

			switch f {
			// --- Validações de CreateProductRequest ---
			case "BrandID", "CategoryID":
				errors[f] = f + " must be a valid UUID."
			case "Name":
				errors[f] = "Name is required and must be between 1 and 255 characters."
			case "Sku":
				errors[f] = "SKU is required and can have a maximum of 50 characters."
			case "BarCodeEan":
				errors[f] = "Barcode EAN can have a maximum of 13 characters."
			case "ShortDescription":
				errors[f] = "Short description can have a maximum of 255 characters."
			case "UnitOfMeasure":
				errors[f] = "Unit of measure is required and can have a maximum of 10 characters."
			case "CostPrice", "SalePrice":
				errors[f] = f + " is required and must be 0 or greater."
			case "PromotionalPrice", "GrossWeight", "NetWeight", "Height", "Width", "Length":
				errors[f] = f + " must be 0 or greater."
			case "Status":
				errors[f] = "Status must be either ACTIVE or INACTIVE."
			case "Stock":
				errors[f] = "Stock object is required."
			case "Fiscal":
				errors[f] = "Fiscal object is required."

			// --- Validações de InventoryProduct (Stock) ---
			case "LocationAisle":
				errors[f] = "Location aisle can have a maximum of 50 characters."
			case "QuantityAvailable", "MinimumStock", "MaximumStock":
				errors[f] = f + " must be 0 or greater."

			// --- Validações de FiscalProduct (Fiscal) ---
			case "NcmCode":
				errors[f] = "NCM code can have a maximum of 8 characters."
			case "CestCode":
				errors[f] = "CEST code can have a maximum of 7 characters."
			case "OriginCode":
				errors[f] = "Origin code must be 0 or greater."
			case "IcmsRate", "PisRate", "CofinsRate", "IpiRate":
				errors[f] = f + " must be between 0 and 100."
			}
		}
	}
	return errors
}