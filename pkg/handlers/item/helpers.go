package item

import (
	errs "checkoutProject/pkg/common/errors"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func addDigitalItemChecks(itemManager ItemManager, log *logrus.Entry, item Item) error {
	isNonDigitalExists, err := itemManager.IsExists(ItemFilter{CategoryIDNot: DIGITAL_ITEM_CATEGORY_ID})
	if err != nil {
		log.WithError(err).Error("error while querying the items")
		return errs.InternalServerErr
	}

	if isNonDigitalExists {
		log.Error("cannot add a digital item if default item exists in cart")
		return fmt.Errorf("cannot add a digital item if default item exists in cart")
	}

	numberOfDigitalItem, err := itemManager.GetTotalItemCount(ItemFilter{CategoryID: DIGITAL_ITEM_CATEGORY_ID})
	if err != nil {
		log.WithError(err).Error("error while finding the number of digital items")
		return errs.InternalServerErr
	}

	if numberOfDigitalItem+item.Quantity > MAX_DIGITAL_ITEMS {
		log.WithError(err).Errorf("error, total number of ditial items cannot be over %d", MAX_DIGITAL_ITEMS)
		return fmt.Errorf("total number of digital items cannot be over %d", MAX_DIGITAL_ITEMS)
	}
	return nil
}

func addDefaultItemChecks(itemManager ItemManager, log *logrus.Entry) error {
	isDigitalItemExists, err := itemManager.IsExists(ItemFilter{CategoryID: DIGITAL_ITEM_CATEGORY_ID})
	if err != nil {
		log.WithError(err).Error("error while querying the items")
		return errs.InternalServerErr
	}

	if isDigitalItemExists {
		log.Error("cannot add a default item if digital item exists in cart")
		return fmt.Errorf("cannot add a default item if digital item exists in cart")
	}

	return nil
}

func addItemPriceChecks(itemManager ItemManager, log *logrus.Entry, item Item) error {
	totalPrice, err := itemManager.GetTotalPrice()
	if err != nil {
		log.WithError(err).Error("error while finding the total price of items")
		return errs.InternalServerErr
	}

	if totalPrice+item.OrderPrice() > MAX_PRICE_OF_CART {
		log.WithError(err).Errorf("total price of cart cannot be over %f", MAX_PRICE_OF_CART)
		return fmt.Errorf("total price of cart cannot be over %.2f", MAX_PRICE_OF_CART)
	}
	return nil
}

func addItemNumberChecks(itemManager ItemManager, log *logrus.Entry, item Item) error {
	numberOfItem, err := itemManager.GetTotalItemCount(ItemFilter{})
	if err != nil {
		log.WithError(err).Error("error while finding the number of items")
		return errs.InternalServerErr
	}

	if item.Quantity+numberOfItem > MAX_DEFAULT_ITEMS {
		log.WithError(err).Errorf("error, total number of items cannot be over %d", MAX_DEFAULT_ITEMS)
		return fmt.Errorf("total number of items cannot be over %d", MAX_DEFAULT_ITEMS)
	}

	numberOfUniqueItem, err := itemManager.GetUniqueItemCount()
	if err != nil {
		log.WithError(err).Error("error while finding the number of unique items")
		return errs.InternalServerErr
	}

	if numberOfUniqueItem >= 10 {
		log.WithError(err).Errorf("error, number of unique items cannot be over %d", MAX_UNIQUE_ITEMS)
		return fmt.Errorf("total number of unique items cannot be over %d", MAX_UNIQUE_ITEMS)
	}
	return nil
}

func addItemIsItemExistsChecks(itemManager ItemManager, log *logrus.Entry, item Item) error {
	isItemExists, err := itemManager.IsExists(ItemFilter{ItemID: item.ItemID})
	if err != nil {
		log.WithError(err).Error("error while querying the item in database")
		return errs.InternalServerErr
	}

	if isItemExists {
		log.Error("error, item with same item id already exists")
		return fmt.Errorf("item with ID %d already exists. Please choose a different item ID", item.ItemID)
	}
	return nil
}

func deleteItemIsItemExistsChecks(itemManager ItemManager, log *logrus.Entry, itemID uint) error {
	isItemExists, err := itemManager.IsExists(ItemFilter{ItemID: itemID})
	if err != nil {
		log.WithError(err).Error("error while querying the item")
		return errs.InternalServerErr
	}

	if !isItemExists {
		log.Error("record not found")
		return errs.RecordNotFoundErr
	}
	return nil
}

func addVasItemIsVasItemExistsInItemChecks(vasItemManager VasItemManager, log *logrus.Entry, vasItemID uint, itemID uint) error {
	isVasItemExistsInItem, err := vasItemManager.IsExistsInItem(ItemVasItemFilter{VasItemID: vasItemID, ItemID: itemID})
	if err != nil {
		log.WithError(err).Error("error while querying the item_vas_item in database")
		return errs.InternalServerErr
	}

	if isVasItemExistsInItem {
		log.Error("error, this item already has this vas-item")
		return fmt.Errorf("item already has this vas-item, cannot add same vas-item multiple times to a single item")
	}
	return nil
}

func addVasItemCategoryAndSellerChecks(log *logrus.Entry, categoryID uint, sellerID uint) error {
	if categoryID != VAS_ITEM_CATEGORY_ID {
		log.Errorf("cannot add vas-item with category id %d", categoryID)
		return fmt.Errorf("cannot add vas-item with category id %d", categoryID)
	}

	if sellerID != VAS_ITEM_SELLER_ID {
		log.Errorf("cannot add vas-item with seller id %d", sellerID)
		return fmt.Errorf("cannot add vas-item with seller id %d", sellerID)
	}
	return nil
}

func addVasItemIsItemExistsAndSuitableChecks(itemManager ItemManager, log *logrus.Entry, itemID uint) (Item, error) {
	item, err := itemManager.Get(ItemFilter{ItemID: itemID})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithError(err).Error("error while querying the item")
		return Item{}, errs.InternalServerErr
	}

	if item.ItemID == 0 {
		log.Error("error, item to add vas-item does not exists")
		return Item{}, fmt.Errorf("cannot add vas-item, item %d does not exist", itemID)
	}

	if !item.isApplicableForVasItems() {
		log.Error("error, item category is not suitable to add vas-items")
		return Item{}, fmt.Errorf("item category is not suitable to add vas-items")
	}
	return item, nil
}

func addVasItemNumberOfVasItemsChecks(itemManager ItemManager, log *logrus.Entry, itemID uint, quantity uint) error {
	numberOfVasItemsInItem, err := itemManager.GetTotalVasItemCount(ItemVasItemFilter{ItemID: itemID})
	if err != nil {
		log.WithError(err).Error("error while finding the number of vas-items in an item")
		return errs.InternalServerErr
	}

	if numberOfVasItemsInItem+quantity > MAX_VAS_ITEM_ON_SINGLE_ITEM {
		log.Errorf("error, cannot add more than %d vas-items to the same item", MAX_VAS_ITEM_ON_SINGLE_ITEM)
		return fmt.Errorf("item %d has already %d vas-items, cannot add more than %d vas-items to the same item", itemID, numberOfVasItemsInItem, MAX_VAS_ITEM_ON_SINGLE_ITEM)
	}
	return nil
}

func addVasItemPriceChecks(itemManager ItemManager, log *logrus.Entry, quantity uint, vasItemPrice float64, itemPrice float64) error {
	totalPrice, err := itemManager.GetTotalPrice()
	if err != nil {
		log.WithError(err).Error("error while finding the total price of the cart")
		return errs.InternalServerErr
	}

	if totalPrice+float64(quantity)*vasItemPrice > MAX_PRICE_OF_CART {
		log.Error("error, vas-items price cannot be more than items price")
		return fmt.Errorf("total price of the cart cannot be ovwer %.2f", MAX_PRICE_OF_CART)
	}

	if itemPrice < vasItemPrice {
		log.Error("error, vas-items price cannot be more than items price")
		return fmt.Errorf("error, sinlge vas-item's price cannot be more than single item's price")
	}
	return nil
}
