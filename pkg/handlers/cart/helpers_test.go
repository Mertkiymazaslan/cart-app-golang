package cart

import (
	errs "checkoutProject/pkg/common/errors"
	"checkoutProject/pkg/common/logger"
	"checkoutProject/pkg/handlers/item"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFindItemsAndVasItems(t *testing.T) {
	log, err := logger.Initialize()
	if err != nil {
		t.Fail()
	}
	mockItemManager := item.NewMockItemManager()
	mockVasItemManager := item.NewMockVasItemManager()

	Convey("TEST getVasItemsOfAnItem fail", t, func() {
		mockItemManager.MFind = func(filter item.ItemFilter) ([]item.Item, error) {
			return []item.Item{
				{ItemID: 1},
			}, nil
		}
		mockVasItemManager.MGetVasItemsOfAnItem = func(filter item.ItemVasItemFilter) ([]item.VasItem, error) {
			return nil, errs.InternalServerErr
		}

		_, err := findItemsAndVasItems(mockItemManager, mockVasItemManager, log.WithFields(logrus.Fields{}))
		So(err, ShouldEqual, errs.InternalServerErr)
	})
	Convey("TEST find items fail", t, func() {
		mockItemManager.MFind = func(filter item.ItemFilter) ([]item.Item, error) {
			return []item.Item{}, errs.InternalServerErr
		}
		mockVasItemManager.MGetVasItemsOfAnItem = func(filter item.ItemVasItemFilter) ([]item.VasItem, error) {
			return []item.VasItem{}, nil
		}

		_, err := findItemsAndVasItems(mockItemManager, mockVasItemManager, log.WithFields(logrus.Fields{}))
		So(err, ShouldEqual, errs.InternalServerErr)
	})
	Convey("TEST success", t, func() {
		mockItemManager.MFind = func(filter item.ItemFilter) ([]item.Item, error) {
			return []item.Item{
				{
					ItemID: 1,
				},
				{
					ItemID: 2,
				},
			}, nil
		}
		mockVasItemManager.MGetVasItemsOfAnItem = func(filter item.ItemVasItemFilter) ([]item.VasItem, error) {
			return []item.VasItem{
				{
					VasItemID: 2,
				},
				{
					VasItemID: 3,
				},
			}, nil
		}

		res, err := findItemsAndVasItems(mockItemManager, mockVasItemManager, log.WithFields(logrus.Fields{}))
		So(err, ShouldBeNil)
		So(res[0].Item.ItemID, ShouldEqual, 1)
		So(res[0].VasItems[0].VasItem.VasItemID, ShouldEqual, 2)
		So(res[0].VasItems[1].VasItem.VasItemID, ShouldEqual, 3)
		So(res[1].VasItems[0].VasItem.VasItemID, ShouldEqual, 2)
		So(res[1].VasItems[1].VasItem.VasItemID, ShouldEqual, 3)
	})
}
