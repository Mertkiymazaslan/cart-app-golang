package item

import (
	errs "checkoutProject/pkg/common/errors"
	"checkoutProject/pkg/common/logger"
	"fmt"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"
	"testing"
)

func TestAddDigitalItemChecks(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}
	mockItemManager := NewMockItemManager()

	Convey("TEST itemManager.isExist fail", t, func() {
		mockItemManager.MIsExists = func(filter ItemFilter) (bool, error) {
			return false, errs.InternalServerErr
		}

		err := addDigitalItemChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{})
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST itemManager.getTotalItemCount fail", t, func() {

		mockItemManager.MIsExists = func(filter ItemFilter) (bool, error) {
			return false, nil
		}

		mockItemManager.MGetTotalItemCount = func(filter ItemFilter) (uint, error) {
			return 0, errs.InternalServerErr
		}

		err := addDigitalItemChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{})
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST item already exists error", t, func() {

		mockItemManager.MIsExists = func(filter ItemFilter) (bool, error) {
			return true, nil
		}

		err := addDigitalItemChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{})
		So(err, ShouldEqual, fmt.Errorf("cannot add a digital item if default item exists in cart"))
	})

	Convey("TEST number of item exceeds limit error", t, func() {

		mockItemManager.MIsExists = func(filter ItemFilter) (bool, error) {
			return false, nil
		}

		mockItemManager.MGetTotalItemCount = func(filter ItemFilter) (uint, error) {
			return 3, nil
		}

		err := addDigitalItemChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{Quantity: 3})
		So(err, ShouldEqual, fmt.Errorf("total number of digital items cannot be over %d", MAX_DIGITAL_ITEMS))
	})

	Convey("TEST succeed and return without error", t, func() {

		mockItemManager.MIsExists = func(filter ItemFilter) (bool, error) {
			return false, nil
		}

		mockItemManager.MGetTotalItemCount = func(filter ItemFilter) (uint, error) {
			return 2, nil
		}

		err := addDigitalItemChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{Quantity: 1})
		So(err, ShouldBeNil)
	})

}

func TestAddDefaultItemChecks(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}
	mockItemManager := NewMockItemManager()

	Convey("TEST itemManager.isExist fail", t, func() {
		mockItemManager.MIsExists = func(filter ItemFilter) (bool, error) {
			return false, errs.InternalServerErr
		}

		err := addDefaultItemChecks(mockItemManager, log.WithFields(logrus.Fields{}))
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST digital items exist error", t, func() {

		mockItemManager.MIsExists = func(filter ItemFilter) (bool, error) {
			return true, nil
		}

		err := addDefaultItemChecks(mockItemManager, log.WithFields(logrus.Fields{}))
		So(err, ShouldEqual, fmt.Errorf("cannot add a default item if digital item exists in cart"))
	})

	Convey("TEST succeed without error", t, func() {

		mockItemManager.MIsExists = func(filter ItemFilter) (bool, error) {
			return false, nil
		}

		err := addDefaultItemChecks(mockItemManager, log.WithFields(logrus.Fields{}))
		So(err, ShouldEqual, nil)
	})
}

func TestAddItemPriceChecks(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}
	mockItemManager := NewMockItemManager()

	Convey("TEST itemManager.getTotalPrice fail", t, func() {
		mockItemManager.MGetTotalPrice = func() (float64, error) {
			return 0, errs.InternalServerErr
		}

		err := addItemPriceChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{})
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST total price exceeds limit error", t, func() {
		mockItemManager.MGetTotalPrice = func() (float64, error) {
			return 300000, nil
		}

		err := addItemPriceChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{Quantity: 2, Price: 100001})
		So(err, ShouldEqual, fmt.Errorf("total price of cart cannot be over %.2f", MAX_PRICE_OF_CART))
	})

	Convey("TEST succeed without error", t, func() {
		mockItemManager.MGetTotalPrice = func() (float64, error) {
			return 300000, nil
		}

		err := addItemPriceChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{Quantity: 2, Price: 100000})
		So(err, ShouldEqual, nil)
	})
}

func TestAddItemNumberNumberChecks(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}
	mockItemManager := NewMockItemManager()

	Convey("TEST itemManager.getTotalItemCount fail", t, func() {
		mockItemManager.MGetTotalItemCount = func(filter ItemFilter) (uint, error) {
			return 0, errs.InternalServerErr
		}

		err := addItemNumberChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{})
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST itemManager.getUniqueItemCount fail", t, func() {
		mockItemManager.MGetTotalItemCount = func(filter ItemFilter) (uint, error) {
			return 25, nil
		}

		mockItemManager.MGetUniqueItemCount = func() (int64, error) {
			return 0, errs.InternalServerErr
		}

		err := addItemNumberChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{})
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST number of item exceeds limit error", t, func() {
		mockItemManager.MGetTotalItemCount = func(filter ItemFilter) (uint, error) {
			return 25, nil
		}

		err := addItemNumberChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{Quantity: 6})
		So(err, ShouldEqual, fmt.Errorf("total number of items cannot be over %d", MAX_DEFAULT_ITEMS))
	})

	Convey("TEST number of unique item exceeds limit error", t, func() {
		mockItemManager.MGetTotalItemCount = func(filter ItemFilter) (uint, error) {
			return 20, nil
		}

		mockItemManager.MGetUniqueItemCount = func() (int64, error) {
			return 10, nil
		}

		err := addItemNumberChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{})
		So(err, ShouldEqual, fmt.Errorf("total number of unique items cannot be over %d", MAX_UNIQUE_ITEMS))
	})

	Convey("TEST succeed without error", t, func() {
		mockItemManager.MGetTotalItemCount = func(filter ItemFilter) (uint, error) {
			return 20, nil
		}

		mockItemManager.MGetUniqueItemCount = func() (int64, error) {
			return 5, nil
		}

		err := addItemNumberChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{})
		So(err, ShouldEqual, nil)
	})
}

func TestAddItemIsItemExistsChecks(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}
	mockItemManager := NewMockItemManager()

	Convey("TEST itemManager.getTotalItemCount fail", t, func() {
		mockItemManager.MIsExists = func(filter ItemFilter) (bool, error) {
			return false, errs.InternalServerErr
		}

		err := addItemIsItemExistsChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{})
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST item already exists error", t, func() {
		mockItemManager.MIsExists = func(filter ItemFilter) (bool, error) {
			return true, nil
		}

		err := addItemIsItemExistsChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{ItemID: 5})
		So(err, ShouldEqual, fmt.Errorf("item with ID 5 already exists. Please choose a different item ID"))
	})

	Convey("TEST succeed withour error", t, func() {
		mockItemManager.MIsExists = func(filter ItemFilter) (bool, error) {
			return false, nil
		}

		err := addItemIsItemExistsChecks(mockItemManager, log.WithFields(logrus.Fields{}), Item{ItemID: 5})
		So(err, ShouldEqual, nil)
	})
}

func TestDeleteItemIsItemExistsChecks(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}
	mockItemManager := NewMockItemManager()

	Convey("TEST itemManager.isExist fail", t, func() {
		mockItemManager.MIsExists = func(filter ItemFilter) (bool, error) {
			return false, errs.InternalServerErr
		}

		err := deleteItemIsItemExistsChecks(mockItemManager, log.WithFields(logrus.Fields{}), 3)
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST item not found error", t, func() {
		mockItemManager.MIsExists = func(filter ItemFilter) (bool, error) {
			return false, nil
		}

		err := deleteItemIsItemExistsChecks(mockItemManager, log.WithFields(logrus.Fields{}), 3)
		So(err, ShouldEqual, errs.RecordNotFoundErr)
	})

	Convey("TEST succeed without error", t, func() {
		mockItemManager.MIsExists = func(filter ItemFilter) (bool, error) {
			return true, nil
		}

		err := deleteItemIsItemExistsChecks(mockItemManager, log.WithFields(logrus.Fields{}), 3)
		So(err, ShouldBeNil)
	})
}

func TestAddVasItemIsVasItemExistsInItemChecks(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}
	mockVasItemManager := NewMockVasItemManager()

	Convey("TEST vasItemManager.isExist fail", t, func() {
		mockVasItemManager.MIsExistsInItem = func(filter ItemVasItemFilter) (bool, error) {
			return false, errs.InternalServerErr
		}

		err := addVasItemIsVasItemExistsInItemChecks(mockVasItemManager, log.WithFields(logrus.Fields{}), 2, 3)
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST vas item exists in item eror", t, func() {
		mockVasItemManager.MIsExistsInItem = func(filter ItemVasItemFilter) (bool, error) {
			return true, nil
		}

		err := addVasItemIsVasItemExistsInItemChecks(mockVasItemManager, log.WithFields(logrus.Fields{}), 2, 3)
		So(err, ShouldEqual, fmt.Errorf("item already has this vas-item, cannot add same vas-item multiple times to a single item"))
	})

	Convey("TEST succeed without error", t, func() {
		mockVasItemManager.MIsExistsInItem = func(filter ItemVasItemFilter) (bool, error) {
			return false, nil
		}

		err := addVasItemIsVasItemExistsInItemChecks(mockVasItemManager, log.WithFields(logrus.Fields{}), 2, 3)
		So(err, ShouldEqual, nil)
	})
}

func TestAddVasItemCategoryAndSellerChecks(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}

	Convey("TEST category_id error", t, func() {
		err := addVasItemCategoryAndSellerChecks(log.WithFields(logrus.Fields{}), 2, 3)
		So(err, ShouldEqual, fmt.Errorf("cannot add vas-item with category id 2"))
	})

	Convey("TEST seller_id error", t, func() {
		err := addVasItemCategoryAndSellerChecks(log.WithFields(logrus.Fields{}), VAS_ITEM_CATEGORY_ID, 3)
		So(err, ShouldEqual, fmt.Errorf("cannot add vas-item with seller id 3"))
	})

	Convey("TEST succeed withour error", t, func() {
		err := addVasItemCategoryAndSellerChecks(log.WithFields(logrus.Fields{}), VAS_ITEM_CATEGORY_ID, VAS_ITEM_SELLER_ID)
		So(err, ShouldEqual, nil)
	})
}

func TestAddVasItemIsItemExistsAndSuitableChecks(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}

	mockItemManager := NewMockItemManager()

	Convey("TEST itemManager.Get fail", t, func() {
		mockItemManager.MGet = func(filter ItemFilter) (Item, error) {
			return Item{}, errs.InternalServerErr
		}

		_, err := addVasItemIsItemExistsAndSuitableChecks(mockItemManager, log.WithFields(logrus.Fields{}), 2)
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST item does not exists error", t, func() {
		mockItemManager.MGet = func(filter ItemFilter) (Item, error) {
			return Item{}, gorm.ErrRecordNotFound
		}

		_, err := addVasItemIsItemExistsAndSuitableChecks(mockItemManager, log.WithFields(logrus.Fields{}), 2)
		So(err, ShouldEqual, fmt.Errorf("cannot add vas-item, item 2 does not exist"))
	})

	Convey("TEST item exists but not suitable to add vas-items error", t, func() {
		mockItemManager.MGet = func(filter ItemFilter) (Item, error) {
			return Item{ItemID: 2, CategoryID: 2}, nil
		}

		_, err := addVasItemIsItemExistsAndSuitableChecks(mockItemManager, log.WithFields(logrus.Fields{}), 2)
		So(err, ShouldEqual, fmt.Errorf("item category is not suitable to add vas-items"))
	})

	Convey("TEST succeed without error", t, func() {
		mockItemManager.MGet = func(filter ItemFilter) (Item, error) {
			return Item{ItemID: 2, CategoryID: FURNITIRE_CATEGORY_ID}, nil
		}

		_, err := addVasItemIsItemExistsAndSuitableChecks(mockItemManager, log.WithFields(logrus.Fields{}), 2)
		So(err, ShouldEqual, nil)
	})
}

func TestAddVasItemNumberOfVasItemsChecks(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}

	mockItemManager := NewMockItemManager()

	Convey("TEST itemManager.GetTotalVasItemCount fail", t, func() {
		mockItemManager.MGetTotalVasItemCount = func(filter ItemVasItemFilter) (uint, error) {
			return 0, errs.InternalServerErr
		}

		err := addVasItemNumberOfVasItemsChecks(mockItemManager, log.WithFields(logrus.Fields{}), 2, 2)
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST number of vas items exceeds limit error", t, func() {
		mockItemManager.MGetTotalVasItemCount = func(filter ItemVasItemFilter) (uint, error) {
			return 2, nil
		}

		err := addVasItemNumberOfVasItemsChecks(mockItemManager, log.WithFields(logrus.Fields{}), 2, 2)
		So(err, ShouldEqual, fmt.Errorf("item 2 has already 2 vas-items, cannot add more than %d vas-items to the same item", MAX_VAS_ITEM_ON_SINGLE_ITEM))
	})

	Convey("TEST succeed without error", t, func() {
		mockItemManager.MGetTotalVasItemCount = func(filter ItemVasItemFilter) (uint, error) {
			return 2, nil
		}

		err := addVasItemNumberOfVasItemsChecks(mockItemManager, log.WithFields(logrus.Fields{}), 2, 1)
		So(err, ShouldEqual, nil)
	})
}

func TestAddVasItemPriceChecks(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}

	mockItemManager := NewMockItemManager()

	Convey("TEST itemManager.Get fail", t, func() {
		mockItemManager.MGetTotalPrice = func() (float64, error) {
			return 0, errs.InternalServerErr
		}

		err := addVasItemPriceChecks(mockItemManager, log.WithFields(logrus.Fields{}), 2, 300, 500)
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST cart total price exceeds limit error", t, func() {
		mockItemManager.MGetTotalPrice = func() (float64, error) {
			return 400000, nil
		}

		err := addVasItemPriceChecks(mockItemManager, log.WithFields(logrus.Fields{}), 2, 150000, 160000)
		So(err, ShouldEqual, fmt.Errorf("total price of the cart cannot be ovwer %.2f", MAX_PRICE_OF_CART))
	})

	Convey("TEST cart total price exceeds limit error", t, func() {
		mockItemManager.MGetTotalPrice = func() (float64, error) {
			return 400000, nil
		}

		err := addVasItemPriceChecks(mockItemManager, log.WithFields(logrus.Fields{}), 2, 150000, 160000)
		So(err, ShouldEqual, fmt.Errorf("total price of the cart cannot be ovwer %.2f", MAX_PRICE_OF_CART))
	})

	Convey("TEST vas items price bigger than items price error", t, func() {
		mockItemManager.MGetTotalPrice = func() (float64, error) {
			return 100000, nil
		}

		err := addVasItemPriceChecks(mockItemManager, log.WithFields(logrus.Fields{}), 2, 10, 5)
		So(err, ShouldEqual, fmt.Errorf("error, sinlge vas-item's price cannot be more than single item's price"))
	})

	Convey("TEST succeed without error", t, func() {
		mockItemManager.MGetTotalPrice = func() (float64, error) {
			return 100000, nil
		}

		err := addVasItemPriceChecks(mockItemManager, log.WithFields(logrus.Fields{}), 2, 4, 5)
		So(err, ShouldBeNil)
	})
}
