package item

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ItemID     uint
	CategoryID uint
	SellerID   uint
	Price      float64
	Quantity   uint
}

func (item Item) isDigitalItem() bool {
	return item.CategoryID == DIGITAL_ITEM_CATEGORY_ID
}

func (item Item) isDefaultItem() bool {
	return item.CategoryID != DIGITAL_ITEM_CATEGORY_ID
}

func (item Item) OrderPrice() float64 {
	return item.Price * float64(item.Quantity)
}

func (item Item) isApplicableForVasItems() bool {
	applicableCategoryIds := []uint{FURNITIRE_CATEGORY_ID, ELECTRONIC_CATEGORY_ID}

	for _, categoryID := range applicableCategoryIds {
		if item.CategoryID == categoryID {
			return true
		}
	}
	return false
}

type VasItem struct {
	gorm.Model
	VasItemID  uint
	CategoryID uint
	SellerID   uint
	Price      float64
	Quantity   uint
}

type ItemVasItem struct {
	gorm.Model
	ItemID    uint
	VasItemID uint
}
