package brands_handler

type CreateBrandRequest struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
}