package cart

import (
	errs "checkoutProject/pkg/common/errors"
	"checkoutProject/pkg/handlers/item"
	"github.com/sirupsen/logrus"
)

func findItemsAndVasItems(itemManager item.ItemManager, vasItemManager item.VasItemManager, log *logrus.Entry) ([]item.ItemSerializer, error) {
	items, err := itemManager.Find(item.ItemFilter{})
	if err != nil {
		log.WithError(err).Error("error while querying the items")
		return nil, errs.InternalServerErr
	}

	var itemsToDisplay []item.ItemSerializer

	for _, itm := range items {
		vasItems, err := vasItemManager.GetVasItemsOfAnItem(item.ItemVasItemFilter{ItemID: itm.ItemID})
		if err != nil {
			log.WithError(err).Error("error while finding vas-items of an item")
			return nil, errs.InternalServerErr
		}

		var vasItemsToDisplay []item.VasItemSerializer
		for _, vasitm := range vasItems {
			vasItemsToDisplay = append(vasItemsToDisplay, item.VasItemSerializer{VasItem: vasitm})
		}

		itemsToDisplay = append(itemsToDisplay, item.ItemSerializer{
			Item:     itm,
			VasItems: vasItemsToDisplay,
		})
	}

	return itemsToDisplay, nil
}
