package cart

import (
	errs "checkoutProject/pkg/common/errors"
	"checkoutProject/pkg/common/logger"
	"checkoutProject/pkg/handlers/item"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestApplyPromotion(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}

	mockItemManager := item.NewMockItemManager()

	Convey("TEST itemManager.find fail", t, func() {
		mockItemManager.MFind = func(filter item.ItemFilter) ([]item.Item, error) {
			return nil, errs.InternalServerErr
		}
		mockItemManager.MAreAllItemsFromSameSeller = func() (bool, error) {
			return true, nil
		}

		_, _, err := ApplyPromotion(500, mockItemManager, log.WithFields(logrus.Fields{}))
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST itemManager.areAllItemsFromSameSeller fail", t, func() {
		mockItemManager.MFind = func(filter item.ItemFilter) ([]item.Item, error) {
			return []item.Item{}, nil
		}
		mockItemManager.MAreAllItemsFromSameSeller = func() (bool, error) {
			return false, errs.InternalServerErr
		}

		_, _, err := ApplyPromotion(500, mockItemManager, log.WithFields(logrus.Fields{}))
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST success and choose same seller promotion", t, func() {
		mockItemManager.MFind = func(filter item.ItemFilter) ([]item.Item, error) {
			return []item.Item{}, nil
		}
		mockItemManager.MAreAllItemsFromSameSeller = func() (bool, error) {
			return true, nil
		}

		discount, promID, err := ApplyPromotion(4000, mockItemManager, log.WithFields(logrus.Fields{}))
		So(err, ShouldEqual, nil)
		So(discount, ShouldEqual, 400)
		So(promID, ShouldEqual, SAME_SELLER_PROMOTION_ID)
	})

	Convey("TEST success and choose category promotion", t, func() {
		mockItemManager.MFind = func(filter item.ItemFilter) ([]item.Item, error) {
			return []item.Item{
				{
					ItemID:     2,
					CategoryID: CATEGORY_PROMOTION_APPLICABLE_CAT_ID,
					Quantity:   2,
					Price:      22000,
					SellerID:   2,
				},
			}, nil
		}
		mockItemManager.MAreAllItemsFromSameSeller = func() (bool, error) {
			return false, nil
		}

		discount, promID, err := ApplyPromotion(44000, mockItemManager, log.WithFields(logrus.Fields{}))
		So(err, ShouldEqual, nil)
		So(discount, ShouldEqual, 2200)
		So(promID, ShouldEqual, CATEGORY_PROMOTION_ID)
	})

	Convey("TEST success and choose total price promotion", t, func() {
		mockItemManager.MFind = func(filter item.ItemFilter) ([]item.Item, error) {
			return []item.Item{
				{
					ItemID:     2,
					CategoryID: CATEGORY_PROMOTION_APPLICABLE_CAT_ID,
					Quantity:   2,
					Price:      50,
					SellerID:   2,
				},
			}, nil
		}
		mockItemManager.MAreAllItemsFromSameSeller = func() (bool, error) {
			return true, nil
		}

		discount, promID, err := ApplyPromotion(360, mockItemManager, log.WithFields(logrus.Fields{}))
		So(err, ShouldEqual, nil)
		So(discount, ShouldEqual, 250)
		So(promID, ShouldEqual, TOTAL_PRICE_PROMOTION_ID)
	})
}

func TestGetSameSellerPromotionDiscount(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}

	mockItemManager := item.NewMockItemManager()

	Convey("TEST areAllItemsFromSameSeller fail", t, func() {
		mockItemManager.MAreAllItemsFromSameSeller = func() (bool, error) {
			return false, errs.InternalServerErr
		}

		_, err := getSameSellerPromotionDiscount(mockItemManager, log.WithFields(logrus.Fields{}), 500)
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST areAllItemsFromSameSeller returns false", t, func() {
		mockItemManager.MAreAllItemsFromSameSeller = func() (bool, error) {
			return false, nil
		}

		discount, err := getSameSellerPromotionDiscount(mockItemManager, log.WithFields(logrus.Fields{}), 500)
		So(err, ShouldBeNil)
		So(discount, ShouldEqual, 0)
	})

	Convey("TEST areAllItemsFromSameSeller returns true", t, func() {
		mockItemManager.MAreAllItemsFromSameSeller = func() (bool, error) {
			return true, nil
		}

		discount, err := getSameSellerPromotionDiscount(mockItemManager, log.WithFields(logrus.Fields{}), 500)
		So(err, ShouldBeNil)
		So(discount, ShouldEqual, 50)
	})

}

func TestGetCategoryPromotionDiscount(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}

	mockItemManager := item.NewMockItemManager()

	Convey("TEST itemManager.find fail", t, func() {
		mockItemManager.MFind = func(filter item.ItemFilter) ([]item.Item, error) {
			return nil, errs.InternalServerErr
		}

		_, err := getCategoryPromotionDiscount(mockItemManager, log.WithFields(logrus.Fields{}))
		So(err, ShouldEqual, errs.InternalServerErr)
	})

	Convey("TEST itemManager.find return no item", t, func() {
		mockItemManager.MFind = func(filter item.ItemFilter) ([]item.Item, error) {
			return []item.Item{}, nil
		}

		discount, err := getCategoryPromotionDiscount(mockItemManager, log.WithFields(logrus.Fields{}))
		So(err, ShouldBeNil)
		So(discount, ShouldEqual, 0)
	})

	Convey("TEST itemManager.find return some items", t, func() {
		mockItemManager.MFind = func(filter item.ItemFilter) ([]item.Item, error) {
			return []item.Item{
				{
					ItemID:     2,
					CategoryID: CATEGORY_PROMOTION_APPLICABLE_CAT_ID,
					Quantity:   3,
					Price:      200,
				},
				{
					ItemID:     2,
					CategoryID: CATEGORY_PROMOTION_APPLICABLE_CAT_ID,
					Quantity:   7,
					Price:      150,
				},
			}, nil
		}

		discount, err := getCategoryPromotionDiscount(mockItemManager, log.WithFields(logrus.Fields{}))
		So(err, ShouldBeNil)
		So(discount, ShouldEqual, 1650*CATEGORY_PROMOTION_PERCENTAGE)
	})
}

func TestGetTotalPricePromotionDiscount(t *testing.T) {

	Convey("TEST between 0 - 5000", t, func() {
		discount := getTotalPricePromotionDiscount(2000)
		So(discount, ShouldEqual, 250)

		discount = getTotalPricePromotionDiscount(200)
		So(discount, ShouldEqual, 200)
	})

	Convey("TEST between 5000 - 10000", t, func() {
		discount := getTotalPricePromotionDiscount(5000)
		So(discount, ShouldEqual, 500)

		getTotalPricePromotionDiscount(8000)
		So(discount, ShouldEqual, 500)
	})

	Convey("TEST between 10000 - 50000", t, func() {
		discount := getTotalPricePromotionDiscount(10000)
		So(discount, ShouldEqual, 1000)

		discount = getTotalPricePromotionDiscount(45000)
		So(discount, ShouldEqual, 1000)
	})

	Convey("TEST 50000+", t, func() {
		discount := getTotalPricePromotionDiscount(50000)
		So(discount, ShouldEqual, 2000)

		discount = getTotalPricePromotionDiscount(264000)
		So(discount, ShouldEqual, 2000)
	})

}
