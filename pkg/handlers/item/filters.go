package item

import "gorm.io/gorm"

type ItemFilter struct {
	ID            uint
	ItemID        uint
	CategoryID    uint
	CategoryIDNot uint
	SellerID      uint
	Price         float64
	Quantity      uint
}

func (f ItemFilter) ToQuery(q *gorm.DB) *gorm.DB {
	q = q.Where(Item{
		Model: gorm.Model{ID: f.ID},
	})

	if f.ItemID != 0 {
		q = q.Where("items.item_id = ?", f.ItemID)
	}

	if f.CategoryID != 0 {
		q = q.Where("items.category_id = ?", f.CategoryID)
	}

	if f.CategoryIDNot != 0 {
		q = q.Where("items.category_id <> ?", f.CategoryIDNot)
	}

	if f.SellerID != 0 {
		q = q.Where("items.seller_id = ?", f.SellerID)
	}

	if f.Price != 0 {
		q = q.Where("items.price = ?", f.Price)
	}

	if f.Quantity != 0 {
		q = q.Where("items.quantity = ?", f.Quantity)
	}

	return q
}

type VasItemFilter struct {
	ID            uint
	VasItemID     uint
	CategoryID    uint
	CategoryIDNot uint
	SellerID      uint
	Price         float64
	Quantity      uint
}

func (f VasItemFilter) ToQuery(q *gorm.DB) *gorm.DB {
	q = q.Where(VasItem{
		Model: gorm.Model{ID: f.ID},
	})

	if f.VasItemID != 0 {
		q = q.Where("vas_items.vas_item_id = ?", f.VasItemID)
	}

	if f.CategoryID != 0 {
		q = q.Where("vas_items.category_id = ?", f.CategoryID)
	}

	if f.CategoryIDNot != 0 {
		q = q.Where("vas_items.category_id <> ?", f.CategoryIDNot)
	}

	if f.SellerID != 0 {
		q = q.Where("vas_items.seller_id = ?", f.SellerID)
	}

	if f.Price != 0 {
		q = q.Where("vas_items.price = ?", f.Price)
	}

	if f.Quantity != 0 {
		q = q.Where("vas_items.quantity = ?", f.Quantity)
	}

	return q
}

type ItemVasItemFilter struct {
	ID        uint
	VasItemID uint
	ItemID    uint
}

func (f ItemVasItemFilter) ToQuery(q *gorm.DB) *gorm.DB {
	q = q.Where(ItemVasItem{
		Model: gorm.Model{ID: f.ID},
	})

	if f.VasItemID != 0 {
		q = q.Where("item_vas_items.vas_item_id = ?", f.VasItemID)
	}

	if f.ItemID != 0 {
		q = q.Where("item_vas_items.item_id = ?", f.ItemID)
	}

	return q
}
