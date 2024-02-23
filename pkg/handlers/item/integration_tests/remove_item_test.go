package integration_tests

import (
	"checkoutProject/pkg/bootstrap"
	"checkoutProject/pkg/common/apiresponse"
	"checkoutProject/pkg/common/testhelper"
	"encoding/json"
	"fmt"
	"github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

type removeItemTest struct {
	Name                    string
	ItemID                  uint
	ExpectedResponseResult  bool
	ExpectedResponseMessage string
	WantCode                int
}

func TestRemoveItem(t *testing.T) {
	tests := []removeItemTest{
		{
			Name:                    "Server should return 400 if item does not exists in cart.",
			ItemID:                  10,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: "record not found",
			WantCode:                http.StatusNotFound,
		},
		{
			Name:                    "Server should return 200 if item successfully deleted from cart.",
			ItemID:                  9,
			ExpectedResponseResult:  true,
			ExpectedResponseMessage: "item removed successfully",
			WantCode:                http.StatusOK,
		},
	}

	db, err := TestDB.DB()
	if err != nil {
		t.Fatalf("error while getting db instance: %v", err)
		return
	}

	testhelper.LoadFixtures(testhelper.DefaultItemsFixturePath, t, db)

	r := bootstrap.SetupRouter()

	for _, test := range tests {
		runRemoveItemTestCase(t, test, r)
	}
}

func runRemoveItemTestCase(t *testing.T, tt removeItemTest, r *gin.Engine) {
	var response gofight.HTTPResponse

	t.Run(tt.Name, func(t *testing.T) {
		gofight.New().
			DELETE(fmt.Sprintf("/api/cart/items/%d", tt.ItemID)).
			Run(r, func(r gofight.HTTPResponse, request gofight.HTTPRequest) {
				response = r
			})

		Convey("When client sends a request to delete the given item", t, func() {
			Convey(fmt.Sprintf("Then server should return %d code", tt.WantCode), func() {
				So(response.Code, ShouldEqual, tt.WantCode)
			})

			var res apiresponse.GenericResponse
			err := json.Unmarshal(response.Body.Bytes(), &res)
			So(err, ShouldBeNil)

			Convey(fmt.Sprintf("Then response should have Result field equal to expected result value"), func() {
				So(res.Result, ShouldEqual, tt.ExpectedResponseResult)
			})

			Convey(fmt.Sprintf("Then response should have Message field equal to expected result message"), func() {
				So(res.Message, ShouldEqual, tt.ExpectedResponseMessage)
			})
			if response.Code == http.StatusOK {
				Convey(fmt.Sprintf("Then item must be deleted in test db if operation is successful"), func() {
					var count int64
					err := TestDB.Table("items").Where("items.item_id = ?", tt.ItemID).Where("items.deleted_at IS NULL").Count(&count).Error
					So(err, ShouldBeNil)
					So(count, ShouldEqual, 0)
				})

				Convey(fmt.Sprintf("Then vas-items related with this items must be deleted if operation is successfull"), func() {
					var count int64
					err := TestDB.Table("item_vas_items").Where("item_vas_items.item_id = ?", tt.ItemID).Where("item_vas_items.deleted_at IS NULL").Count(&count).Error
					So(err, ShouldBeNil)
					So(count, ShouldEqual, 0)
				})
			}
		})
	})
}
