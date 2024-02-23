package integration_tests

import (
	"checkoutProject/pkg/bootstrap"
	"checkoutProject/pkg/common/apiresponse"
	"checkoutProject/pkg/common/testhelper"
	itm "checkoutProject/pkg/handlers/item"
	"encoding/json"
	"fmt"
	"github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

type addItemTest struct {
	Name                    string
	ItemID                  uint
	CategoryID              uint
	SellerID                uint
	Price                   float64
	Quantity                uint
	ExpectedResponseResult  bool
	ExpectedResponseMessage string
	WantCode                int
}

func TestAddItemForCartWithDefaultItems(t *testing.T) {
	tests := []addItemTest{
		{
			Name:                    "server should return 400 if item already exists",
			ItemID:                  1,
			CategoryID:              1,
			SellerID:                1,
			Price:                   16.34,
			Quantity:                3,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: "item with ID 1 already exists. Please choose a different item ID",
			WantCode:                http.StatusBadRequest,
		},
		{
			Name:                    "server should return 400 if client tries to add digital item to cart with default items",
			ItemID:                  100,
			CategoryID:              itm.DIGITAL_ITEM_CATEGORY_ID,
			SellerID:                1,
			Price:                   10000,
			Quantity:                3,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: "cannot add a digital item if default item exists in cart",
			WantCode:                http.StatusBadRequest,
		},
		{
			Name:                    "server should return 400 if client tries to add items that make the carts total price bigger than the limit",
			ItemID:                  101,
			CategoryID:              10,
			SellerID:                1,
			Price:                   25000,
			Quantity:                2,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: fmt.Sprintf("total price of cart cannot be over %.2f", itm.MAX_PRICE_OF_CART),
			WantCode:                http.StatusBadRequest,
		},
		{
			Name:                    "server should return 400 if client tries to add more than 30 items in cart",
			ItemID:                  102,
			CategoryID:              1,
			SellerID:                1,
			Price:                   50.7,
			Quantity:                8,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: fmt.Sprintf("total number of items cannot be over %d", itm.MAX_DEFAULT_ITEMS),
			WantCode:                http.StatusBadRequest,
		},
		{
			Name:                    "server should return 201 if item added successfully",
			ItemID:                  103,
			CategoryID:              1,
			SellerID:                1,
			Price:                   10000,
			Quantity:                1,
			ExpectedResponseResult:  true,
			ExpectedResponseMessage: "item added successfully",
			WantCode:                http.StatusCreated,
		},
		{
			Name:                    "server should return 400 if client tries to add more than 10 unique items in cart",
			ItemID:                  104,
			CategoryID:              1,
			SellerID:                1,
			Price:                   10000,
			Quantity:                1,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: fmt.Sprintf("total number of unique items cannot be over %d", itm.MAX_UNIQUE_ITEMS),
			WantCode:                http.StatusBadRequest,
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
		runAddItemTestCase(t, test, r)
	}
}

func TestAddItemForCartWithDigitalItems(t *testing.T) {
	tests := []addItemTest{
		{
			Name:                    "server should return 400 if client try to add default item to the cart with digital item(s).",
			ItemID:                  10,
			CategoryID:              1,
			SellerID:                1,
			Price:                   16.34,
			Quantity:                1,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: "cannot add a default item if digital item exists in cart",
			WantCode:                http.StatusBadRequest,
		},
		{
			Name:                    "server should return 400 if client try to add more than 5 digital items to cart.",
			ItemID:                  20,
			CategoryID:              itm.DIGITAL_ITEM_CATEGORY_ID,
			SellerID:                1,
			Price:                   16.34,
			Quantity:                2,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: fmt.Sprintf("total number of digital items cannot be over %d", itm.MAX_DIGITAL_ITEMS),
			WantCode:                http.StatusBadRequest,
		},
		{
			Name:                    "server should return 201 if item added successfully.",
			ItemID:                  11,
			CategoryID:              itm.DIGITAL_ITEM_CATEGORY_ID,
			SellerID:                56,
			Price:                   10000,
			Quantity:                1,
			ExpectedResponseResult:  true,
			ExpectedResponseMessage: fmt.Sprintf("item added successfully"),
			WantCode:                http.StatusCreated,
		},

		//I added some detailed/rare test cases here to avoid extending the code. (not the best practices)
		{
			Name:                    "server should return 400 and give details about missing fields.",
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: "{\"CategoryID\":\"This field is required\",\"ItemID\":\"This field is required\",\"Price\":\"This field is required\",\"Quantity\":\"This field is required\",\"SellerID\":\"This field is required\"}",
			WantCode:                http.StatusBadRequest,
		},
		{
			Name:                    "server should return 400 and give details about failed min-max binding checks",
			ItemID:                  15,
			CategoryID:              itm.DIGITAL_ITEM_CATEGORY_ID,
			SellerID:                56,
			Price:                   10000,
			Quantity:                63,
			ExpectedResponseResult:  false,
			ExpectedResponseMessage: "{\"Quantity\":\"This fields maximum value is 10\"}",
			WantCode:                http.StatusBadRequest,
		},
	}

	db, err := TestDB.DB()
	if err != nil {
		t.Fatalf("error while getting db instance: %v", err)
		return
	}

	testhelper.LoadFixtures(testhelper.DigitalItemFixturesPath, t, db)

	r := bootstrap.SetupRouter()

	for _, test := range tests {
		runAddItemTestCase(t, test, r)
	}
}

func runAddItemTestCase(t *testing.T, tt addItemTest, r *gin.Engine) {
	var response gofight.HTTPResponse

	t.Run(tt.Name, func(t *testing.T) {
		gofight.New().
			POST("/api/cart/items").
			SetJSON(gofight.D{
				"item_id":     tt.ItemID,
				"category_id": tt.CategoryID,
				"seller_id":   tt.SellerID,
				"price":       tt.Price,
				"quantity":    tt.Quantity,
			}).
			Run(r, func(r gofight.HTTPResponse, request gofight.HTTPRequest) {
				response = r
			})

		Convey("When client sends a request to create a new item", t, func() {
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
				Convey(fmt.Sprintf("Then item must be created in test db if operation is successful"), func() {

					var count int64
					err := TestDB.Table("items").Where("items.item_id = ?", tt.ItemID).Count(&count).Error
					So(err, ShouldBeNil)
					So(count, ShouldEqual, 1)
				})
			}
		})
	})
}
