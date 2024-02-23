package cart

import (
	"checkoutProject/pkg/common/apiresponse"
	db "checkoutProject/pkg/common/database"
	errs "checkoutProject/pkg/common/errors"
	"checkoutProject/pkg/common/logger"
	"checkoutProject/pkg/handlers/item"
	"github.com/sirupsen/logrus"
)

type CartController interface {
	DisplayCart() (apiresponse.Responder, error)
	ResetCart() (apiresponse.Responder, error)
}

type cartController struct {
	itemManager    item.ItemManager
	vasItemManager item.VasItemManager
}

func NewCartController(itemManager item.ItemManager, vasItemManager item.VasItemManager) CartController {
	return cartController{
		itemManager:    itemManager,
		vasItemManager: vasItemManager,
	}
}

func NewDefaultCartController() CartController {
	return NewCartController(item.NewDefaultItemManager(), item.NewDefaultVasItemManager())
}

func (c cartController) formattedLogger(l logrus.FieldLogger) *logrus.Entry {
	return l.WithFields(logrus.Fields{"api_version": "1", "controller": "cart"})
}

func (c cartController) DisplayCart() (apiresponse.Responder, error) {
	log := c.formattedLogger(logger.GetInstance()).WithFields(logrus.Fields{
		"location": "Display Cart",
	})

	itemManager := c.itemManager
	vasItemManager := c.vasItemManager

	itemsToDisplay, err := findItemsAndVasItems(itemManager, vasItemManager, log)
	if err != nil {
		return nil, err
	}

	totalPrice, err := itemManager.GetTotalPrice()
	if err != nil {
		log.WithError(err).Error("error while finding total price")
		return nil, errs.InternalServerErr
	}

	discount, promotionID, err := ApplyPromotion(totalPrice, itemManager, log)

	newPrice := totalPrice - discount

	resp := CartSerializer{Result: true, Message: CartMessageSerializer{
		Items:              itemsToDisplay,
		TotalPrice:         newPrice,
		AppliedPromotionID: promotionID,
		TotalDiscount:      discount,
	}}

	return resp, nil
}

func (c cartController) ResetCart() (apiresponse.Responder, error) {
	log := c.formattedLogger(logger.GetInstance()).WithFields(logrus.Fields{
		"location": "Reset Cart",
	})

	tx := db.NewTransaction()
	defer func() {
		if err := db.RollbackTransaction(tx); err != nil {
			log.WithError(err).Error("error while rolling back transaction")
		}
	}()

	itemManager := c.itemManager.WithTx(tx)
	vasItemManager := c.vasItemManager.WithTx(tx)

	err := vasItemManager.DeleteAllItemVasItems()
	if err != nil {
		log.WithError(err).Error("error while deleting the item_vas_items")
		return nil, errs.InternalServerErr
	}

	err = itemManager.DeleteAllItems()
	if err != nil {
		log.WithError(err).Error("error while deleting the items")
		return nil, errs.InternalServerErr
	}

	err = vasItemManager.DeleteAllVasItems()
	if err != nil {
		log.WithError(err).Error("error while deleting the vas_items")
		return nil, errs.InternalServerErr
	}

	if err = db.CommitTransaction(tx); err != nil {
		log.WithError(err).Error("error while committing the transaction")
		return nil, errs.InternalServerErr
	}

	return apiresponse.GenericResponseSerializer{Result: true, Message: "cart emptied successfully"}, nil
}
