package category_handler

type CreateCategoriesRequest struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
}

type CreateCategoryResponse struct {
	ID       string
	ParentID string
	Name     string
}
