package product

type Product struct {
	ID               string  `json:"id"`
	BrandID          string  `json:"brand_id"`
	CategoryID       string  `json:"category_id"`
	Name             string  `json:"name"`
	Sku              string  `json:"sku"`
	BarcodeEan13     string  `json:"barcode_ean13"`
	ShortDescription string  `json:"short_description"`
	UnitOfMeasure    string  `json:"unit_of_measure"`
	CostPrice        float64 `json:"cost_price"`
	SalePrice        float64 `json:"sale_price"`
	PromotionalPrice float64 `json:"promotional_price"`
	GrossWeight      float64 `json:"gross_weight"`
	NetWeight        float64 `json:"net_weight"`
	Height           float64 `json:"height"`
	Width            int     `json:"width"`
	Length           int     `json:"length"`
	Status           string  `json:"status"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	Stock            struct {
		ID                string `json:"id"`
		ProductID         string `json:"product_id"`
		LocationAisle     string `json:"location_aisle"`
		QuantityAvailable int    `json:"quantity_available"`
		MinimumStock      int    `json:"minimum_stock"`
		MaximumStock      int    `json:"maximum_stock"`
		UpdatedAt         string `json:"updated_at"`
	} `json:"stock"`
	Fiscal struct {
		ID         string  `json:"id"`
		ProductID  string  `json:"product_id"`
		NcmCode    string  `json:"ncm_code"`
		CestCode   string  `json:"cest_code"`
		IcmsRate   int     `json:"icms_rate"`
		PisRate    float64 `json:"pis_rate"`
		CofinsRate float64 `json:"cofins_rate"`
		IpiRate    int     `json:"ipi_rate"`
		UpdatedAt  string  `json:"updated_at"`
	} `json:"fiscal"`
}
