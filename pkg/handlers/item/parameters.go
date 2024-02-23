package item

type ItemUriParams struct {
	ItemID uint `uri:"item_id" binding:"required"`
}

type AddItemParams struct {
	ItemID     uint    `json:"item_id" binding:"required"`
	CategoryID uint    `json:"category_id" binding:"required"`
	SellerID   uint    `json:"seller_id" binding:"required"`
	Price      float64 `json:"price" binding:"required,min=1,max=500000"`
	Quantity   uint    `json:"quantity" binding:"required,min=1,max=10"`
}

type AddVasItemParams struct {
	ItemUriParams
	VasItemID  uint    `json:"vas_item_id" binding:"required"`
	CategoryID uint    `json:"category_id" binding:"required"`
	SellerID   uint    `json:"seller_id" binding:"required"`
	Price      float64 `json:"price" binding:"required,min=1,max=500000"`
	Quantity   uint    `json:"quantity" binding:"required,min=1,max=3"`
}

type RemoveItemParams struct {
	ItemUriParams
}
