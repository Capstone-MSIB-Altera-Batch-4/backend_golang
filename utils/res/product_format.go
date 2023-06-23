package res

import "point-of-sale/app/model"

type SetProductsFormat struct {
	ID         int    `json:"id"`
	ProductsID string `json:"products_id"`
	Image      string `json:"image_url"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	Unit       string `json:"unit"`
	Quantity   int    `json:"quantity"`
	Price      int    `json:"price"`
}

func TransformAdminProduct(products model.Product) SetProductsFormat {
	return SetProductsFormat{
		ID:         products.ID,
		ProductsID: products.ProductID,
		Image:      products.Image,
		Name:       products.Name,
		Category:   products.Category.Name,
		Unit:       products.Unit,
		Quantity:   products.Quantity,
		Price:      products.Price,
	}
}

func TransformAdminProducts(products []model.Product) []SetProductsFormat {
	transformedProducts := make([]SetProductsFormat, len(products))
	for i, p := range products {
		transformedProducts[i] = SetProductsFormat{
			ID:         p.ID,
			ProductsID: p.ProductID,
			Image:      p.Image,
			Name:       p.Name,
			Category:   p.Category.Name,
			Unit:       p.Unit,
			Quantity:   p.Quantity,
			Price:      p.Price,
		}
	}
	return transformedProducts
}
