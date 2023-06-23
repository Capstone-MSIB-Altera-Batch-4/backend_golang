package res

import "point-of-sale/app/model"

type SetCategoryFormat struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

func TransformCategory(request []model.Category) []SetCategoryFormat {
	transformedCategory := make([]SetCategoryFormat, len(request))
	for i, category := range request {
		transformedCategory[i] = SetCategoryFormat{
			ID:   category.ID,
			Name: category.Name,
		}
	}
	return transformedCategory
}
