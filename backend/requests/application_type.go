package requests

type CreateApplicationTypeRequest struct {
	Title       string `json:"title" binding:"required"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateApplicationTypeRequest struct {
	Title       string `json:"title"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
	SortOrder   *int   `json:"sort_order"`
}

type ApplicationTypeFilter struct {
	ActiveOnly bool   `form:"active_only"`
	Search     string `form:"search"`
}
