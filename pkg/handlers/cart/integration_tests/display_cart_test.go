package integration_tests

import (
	"checkoutProject/pkg/bootstrap"
	"checkoutProject/pkg/common/testhelper"
	"checkoutProject/pkg/handlers/cart"
	"checkoutProject/pkg/handlers/item"
	"encoding/json"
	"fmt"
	"github.com/appleboy/gofight/v2"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

type displayCartTest struct {
	Name             string
	ExpectedResponse interface{}
	WantCode         int
}

func TestDisplayCart(t *testing.T) {
	tests := []displayCartTest{
		{
			Name: "Server should return 400 if item does not exists in cart.",
			ExpectedResponse: cart.CartResponse{Result: true, Message: cart.CartMessageResponse{
				Items: []item.ItemResponse{
					{
						ItemID:     1,
						CategoryID: 1001,
						SellerID:   1,
						Price:      20.45,
						Quantity:   1,
						VasItems: []item.VasItemResponse{
							{
								VasItemID:  1,
								CategoryID: 3242,
								SellerID:   5003,
								Price:      50,
								Quantity:   2,
							},
						},
					},
					{
						ItemID:     2,
						CategoryID: 1001,
						SellerID:   1,
						Price:      30.50,
						Quantity:   6,
						VasItems: []item.VasItemResponse{
							{
								VasItemID:  2,
								CategoryID: 3242,
								SellerID:   5003,
								Price:      40.2,
								Quantity:   2,
							},
							{
								VasItemID:  3,
								CategoryID: 3242,
								SellerID:   5003,
								Price:      30.50,
								Quantity:   1,
							},
						},
					},
					{
						ItemID:     3,
						CategoryID: 3004,
						SellerID:   1,
						Price:      3.50,
						Quantity:   1,
						VasItems: []item.VasItemResponse{
							{
								VasItemID:  3,
								CategoryID: 3242,
								SellerID:   5003,
								Price:      30.50,
								Quantity:   1,
							},
						},
					},
					{
						ItemID:     4,
						CategoryID: 1001,
						SellerID:   6,
						Price:      100000,
						Quantity:   2,
						VasItems:   []item.VasItemResponse{},
					},
				},
				TotalPrice:         198448.35,
				AppliedPromotionID: 1232,
				TotalDiscount:      2000,
			}},
			WantCode: http.StatusOK,
		},
	}

	db, err := TestDB.DB()
	if err != nil {
		t.Fatalf("error while getting db instance: %v", err)
		return
	}

	testhelper.LoadFixtures(testhelper.DefaultPath, t, db)

	r := bootstrap.SetupRouter()

	for _, tt := range tests {
		var response gofight.HTTPResponse

		t.Run(tt.Name, func(t *testing.T) {
			gofight.New().
				GET("/api/cart").
				Run(r, func(r gofight.HTTPResponse, request gofight.HTTPRequest) {
					response = r
				})

			Convey("When client sends a request to display the cart info", t, func() {
				Convey(fmt.Sprintf("Then server should return %d code", tt.WantCode), func() {
					So(response.Code, ShouldEqual, tt.WantCode)
				})
				Convey("Then server should return correct response", func() {
					bytes := response.Body.Bytes()
					expectedResponseBytes, err := json.Marshal(tt.ExpectedResponse)
					So(err, ShouldBeNil)

					So(string(bytes), ShouldEqual, string(expectedResponseBytes))
				})
			})
		})
	}
}
