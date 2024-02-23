package item

type ItemResponse struct {
	ItemID     uint              `json:"item_id"`
	CategoryID uint              `json:"category_id"`
	SellerID   uint              `json:"seller_id"`
	Price      float64           `json:"price"`
	Quantity   uint              `json:"quantity"`
	VasItems   []VasItemResponse `json:"vas_items"`
}

type ItemSerializer struct {
	Item     Item
	VasItems []VasItemSerializer
}

func (s ItemSerializer) Response() interface{} {
	vasItems := []VasItemResponse{}
	for _, item := range s.VasItems {
		vasItems = append(vasItems, item.Response().(VasItemResponse))
	}
	return ItemResponse{
		ItemID:     s.Item.ItemID,
		CategoryID: s.Item.CategoryID,
		SellerID:   s.Item.SellerID,
		Price:      s.Item.Price,
		Quantity:   s.Item.Quantity,
		VasItems:   vasItems,
	}
}

type VasItemResponse struct {
	VasItemID  uint    `json:"vas_item_id"`
	CategoryID uint    `json:"category_id"`
	SellerID   uint    `json:"seller_id"`
	Price      float64 `json:"price"`
	Quantity   uint    `json:"quantity"`
}

type VasItemSerializer struct {
	VasItem VasItem
}

func (s VasItemSerializer) Response() interface{} {
	return VasItemResponse{
		VasItemID:  s.VasItem.VasItemID,
		CategoryID: s.VasItem.CategoryID,
		SellerID:   s.VasItem.SellerID,
		Price:      s.VasItem.Price,
		Quantity:   s.VasItem.Quantity,
	}
}
