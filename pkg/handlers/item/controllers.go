package item

import (
	"checkoutProject/pkg/common/apiresponse"
	db "checkoutProject/pkg/common/database"
	errs "checkoutProject/pkg/common/errors"
	"checkoutProject/pkg/common/logger"
	"fmt"
	"github.com/sirupsen/logrus"
)

type ItemController interface {
	AddItem(params AddItemParams) (apiresponse.Responder, error)
	RemoveItem(params RemoveItemParams) (apiresponse.Responder, error)
}

type itemController struct {
	itemManager ItemManager
}

func NewItemController(itemManager ItemManager) ItemController {
	return itemController{
		itemManager: itemManager,
	}
}

func NewDefaultItemController() ItemController {
	return NewItemController(NewDefaultItemManager())
}

func (c itemController) formattedLogger(l logrus.FieldLogger) *logrus.Entry {
	return l.WithFields(logrus.Fields{"api_version": "1", "controller": "item"})
}

func (c itemController) AddItem(params AddItemParams) (apiresponse.Responder, error) {
	log := c.formattedLogger(logger.GetInstance()).WithFields(logrus.Fields{
		"location": "Add Item",
	})

	if params.CategoryID == VAS_ITEM_CATEGORY_ID {
		return nil, fmt.Errorf("cannot add vas-item from this endpoint")
	}

	itemManager := c.itemManager

	item := Item{
		ItemID:     params.ItemID,
		SellerID:   params.SellerID,
		CategoryID: params.CategoryID,
		Price:      params.Price,
		Quantity:   params.Quantity,
	}

	err := addItemIsItemExistsChecks(itemManager, log, item)
	if err != nil {
		return nil, err
	}

	if item.isDigitalItem() {
		err = addDigitalItemChecks(itemManager, log, item)
		if err != nil {
			return nil, err
		}
	}

	if item.isDefaultItem() {
		err = addDefaultItemChecks(itemManager, log)
		if err != nil {
			return nil, err
		}
	}

	err = addItemPriceChecks(itemManager, log, item)
	if err != nil {
		return nil, err
	}

	err = addItemNumberChecks(itemManager, log, item)
	if err != nil {
		return nil, err
	}

	_, err = itemManager.Create(item)
	if err != nil {
		log.WithError(err).Error("error while creating item")
		return nil, errs.InternalServerErr
	}

	return apiresponse.GenericResponseSerializer{Result: true, Message: "item added successfully"}, nil
}

func (c itemController) RemoveItem(params RemoveItemParams) (apiresponse.Responder, error) {

	log := c.formattedLogger(logger.GetInstance()).WithFields(logrus.Fields{
		"location": "Remove Item",
	})

	tx := db.NewTransaction()
	defer func() {
		if err := db.RollbackTransaction(tx); err != nil {
			log.WithError(err).Error("error while rolling back transaction")
		}
	}()

	itemManager := c.itemManager.WithTx(tx)

	err := deleteItemIsItemExistsChecks(itemManager, log, params.ItemID)
	if err != nil {
		return nil, err
	}

	err = itemManager.DeleteVasItemsOfItem(params.ItemID)
	if err != nil {
		log.WithError(err).Error("error while deleting the vas-items of the item")
		return nil, errs.InternalServerErr
	}

	err = itemManager.Delete(ItemFilter{ItemID: params.ItemID})
	if err != nil {
		log.WithError(err).Error("error while deleting the item")
		return nil, errs.InternalServerErr
	}

	if err = db.CommitTransaction(tx); err != nil {
		log.WithError(err).Error("error while committing the transaction")
		return nil, errs.InternalServerErr
	}

	return apiresponse.GenericResponseSerializer{Result: true, Message: "item removed successfully"}, nil
}

// vas-item controller
type VasItemController interface {
	AddVasItem(params AddVasItemParams) (apiresponse.Responder, error)
}

type vasItemController struct {
	vasItemManager VasItemManager
	itemManager    ItemManager
}

func NewVasItemController(vasItemManager VasItemManager, itemManager ItemManager) VasItemController {
	return vasItemController{
		vasItemManager: vasItemManager,
		itemManager:    itemManager,
	}
}

func NewDefaultVasItemController() VasItemController {
	return NewVasItemController(NewDefaultVasItemManager(), NewDefaultItemManager())
}

func (c vasItemController) formattedLogger(l logrus.FieldLogger) *logrus.Entry {
	return l.WithFields(logrus.Fields{"api_version": "1", "controller": "vas_item"})
}

func (c vasItemController) AddVasItem(params AddVasItemParams) (apiresponse.Responder, error) {
	log := c.formattedLogger(logger.GetInstance()).WithFields(logrus.Fields{
		"location": "Add vas item",
	})

	tx := db.NewTransaction()
	defer func() {
		if err := db.RollbackTransaction(tx); err != nil {
			log.WithError(err).Error("error while rolling back transaction")
		}
	}()

	vasItemManager := c.vasItemManager.WithTx(tx)
	itemManager := c.itemManager

	err := addVasItemIsVasItemExistsInItemChecks(vasItemManager, log, params.VasItemID, params.ItemID)
	if err != nil {
		return nil, err
	}

	err = addVasItemCategoryAndSellerChecks(log, params.CategoryID, params.SellerID)
	if err != nil {
		return nil, err
	}

	item, err := addVasItemIsItemExistsAndSuitableChecks(itemManager, log, params.ItemID)
	if err != nil {
		return nil, err
	}

	err = addVasItemNumberOfVasItemsChecks(itemManager, log, params.ItemID, params.Quantity)
	if err != nil {
		return nil, err
	}

	err = addVasItemPriceChecks(itemManager, log, params.Quantity, params.Price, item.Price)
	if err != nil {
		return nil, err
	}

	isVasItemExists, err := vasItemManager.IsExists(VasItemFilter{VasItemID: params.VasItemID})
	if err != nil {
		log.WithError(err).Error("error while querying the vas-item in database")
		return nil, errs.InternalServerErr
	}

	if !isVasItemExists {
		vasItem := VasItem{
			VasItemID:  params.VasItemID,
			SellerID:   params.SellerID,
			CategoryID: params.CategoryID,
			Price:      params.Price,
			Quantity:   params.Quantity,
		}
		_, err := vasItemManager.CreateNewVasItem(vasItem)
		if err != nil {
			log.WithError(err).Error("error while creating new vas-item")
			return nil, errs.InternalServerErr
		}
	}

	itemVasItem := ItemVasItem{VasItemID: params.VasItemID, ItemID: params.ItemID}

	_, err = vasItemManager.CreateItemVasItem(itemVasItem)
	if err != nil {
		log.WithError(err).Error("error while creating the item_vas_item")
		return nil, errs.InternalServerErr
	}

	if err = db.CommitTransaction(tx); err != nil {
		log.WithError(err).Error("error while committing the transaction")
		return nil, errs.InternalServerErr
	}

	return apiresponse.GenericResponseSerializer{Result: true, Message: "vas-item added successfully"}, nil
}
