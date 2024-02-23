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

type addVasItemTest struct {
	Name                    string
	VasItemID               uint
	ItemID                  uint
	CategoryID              uint
	SellerID                uint
	Price                   float64
	Quantity                uint
	ExpectedResponseResult  bool
	ExpectedResponseMessage string
	WantCode                int
}

func TestAddVasItem(t *testing.T) {
	tests := []addVasItemTest{
		{
			Name:                    "server should return 400 if item_vas_item already exists",
			ItemID:                  1,
			VasItemID:               1,
			CategoryID:              3242,
			SellerID:                5003,
			Price:                   16.34,
			Quantity:                1,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: "item already has this vas-item, cannot add same vas-item multiple times to a single item",
			WantCode:                http.StatusBadRequest,
		},
		{
			Name:                    "server should return 400 category_id is not true",
			ItemID:                  2,
			VasItemID:               2,
			CategoryID:              5,
			SellerID:                5003,
			Price:                   16.34,
			Quantity:                1,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: "cannot add vas-item with category id 5",
			WantCode:                http.StatusBadRequest,
		},
		{
			Name:                    "server should return 400 if seller_id is not true",
			ItemID:                  2,
			VasItemID:               2,
			CategoryID:              3242,
			SellerID:                3,
			Price:                   16.34,
			Quantity:                1,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: "cannot add vas-item with seller id 3",
			WantCode:                http.StatusBadRequest,
		},
		{
			Name:                    "server should return 400 if default item does not exists",
			ItemID:                  9999,
			VasItemID:               2,
			CategoryID:              3242,
			SellerID:                5003,
			Price:                   16.34,
			Quantity:                1,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: "cannot add vas-item, item 9999 does not exist",
			WantCode:                http.StatusBadRequest,
		},
		{
			Name:                    "server should return 400 if item category is not suitable to add vas-items",
			ItemID:                  3,
			VasItemID:               2,
			CategoryID:              3242,
			SellerID:                5003,
			Price:                   16.34,
			Quantity:                1,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: "item category is not suitable to add vas-items",
			WantCode:                http.StatusBadRequest,
		},
		{
			Name:                    "server should return 400 if number of vas-items on single item limit exceeded",
			ItemID:                  4,
			VasItemID:               3,
			CategoryID:              3242,
			SellerID:                5003,
			Price:                   16.34,
			Quantity:                2,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: "item 4 has already 2 vas-items, cannot add more than 3 vas-items to the same item",
			WantCode:                http.StatusBadRequest,
		},
		{
			Name:                    "server should return 400 if the price of single vas-item is greater than the default item's price",
			ItemID:                  5,
			VasItemID:               7,
			CategoryID:              3242,
			SellerID:                5003,
			Price:                   16.34,
			Quantity:                2,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: "error, sinlge vas-item's price cannot be more than single item's price",
			WantCode:                http.StatusBadRequest,
		},
		{
			Name:                    "server should return 201 and create new vas-item and item_vas_item",
			ItemID:                  6,
			VasItemID:               10,
			CategoryID:              3242,
			SellerID:                5003,
			Price:                   5,
			Quantity:                2,
			ExpectedResponseResult:  true,
			ExpectedResponseMessage: "vas-item added successfully",
			WantCode:                http.StatusCreated,
		},
		{
			Name:                    "server should return 201 and not create new vas-item since it existed before but create new item_vas_item",
			ItemID:                  5,
			VasItemID:               2,
			CategoryID:              3242,
			SellerID:                5003,
			Price:                   3.6,
			Quantity:                2,
			ExpectedResponseResult:  true,
			ExpectedResponseMessage: "vas-item added successfully",
			WantCode:                http.StatusCreated,
		},
	}

	db, err := TestDB.DB()
	if err != nil {
		t.Fatalf("error while getting db instance: %v", err)
		return
	}

	testhelper.LoadFixtures(testhelper.AddVasItemFixturesPath, t, db)

	r := bootstrap.SetupRouter()

	for _, test := range tests {
		runAddVasItemTestCase(t, test, r)
	}
}

func runAddVasItemTestCase(t *testing.T, tt addVasItemTest, r *gin.Engine) {
	var response gofight.HTTPResponse

	t.Run(tt.Name, func(t *testing.T) {
		gofight.New().
			POST(fmt.Sprintf("/api/cart/items/%d/vas-items", tt.ItemID)).
			SetJSON(gofight.D{
				"vas_item_id": tt.VasItemID,
				"category_id": tt.CategoryID,
				"seller_id":   tt.SellerID,
				"price":       tt.Price,
				"quantity":    tt.Quantity,
			}).
			Run(r, func(r gofight.HTTPResponse, request gofight.HTTPRequest) {
				response = r
			})

		Convey("When client sends a request to create a new vas-item", t, func() {
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

			if response.Code == http.StatusCreated {
				Convey(fmt.Sprintf("Then vas-item must be created or exist before in test db if operation is successful"), func() {
					var count int64
					err := TestDB.Table("vas_items").Where("vas_items.vas_item_id = ?", tt.VasItemID).Count(&count).Error
					So(err, ShouldBeNil)
					So(count, ShouldEqual, 1)
				})

				Convey(fmt.Sprintf("Then item_vas_item must be created in test db if operation is successful"), func() {
					var count int64
					err := TestDB.Table("item_vas_items").Where("item_vas_items.vas_item_id = ? AND item_vas_items.item_id = ?", tt.VasItemID, tt.ItemID).Count(&count).Error
					So(err, ShouldBeNil)
					So(count, ShouldEqual, 1)
				})
			}
		})
	})
}
