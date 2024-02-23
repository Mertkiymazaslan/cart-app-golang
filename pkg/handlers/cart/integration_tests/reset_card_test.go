package integration_tests

import (
	"checkoutProject/pkg/bootstrap"
	"checkoutProject/pkg/common/apiresponse"
	"checkoutProject/pkg/common/testhelper"
	"encoding/json"
	"fmt"
	"github.com/appleboy/gofight/v2"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

type resetCardTest struct {
	Name                    string
	ExpectedResponseResult  bool
	ExpectedResponseMessage string
	WantCode                int
}

func TestResetCard(t *testing.T) {
	tests := []resetCardTest{
		{
			Name:                    "Server should return 200 and reset the cart",
			ExpectedResponseResult:  true,
			ExpectedResponseMessage: "cart emptied successfully",
			WantCode:                http.StatusOK,
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
				DELETE("/api/cart/reset").
				Run(r, func(r gofight.HTTPResponse, request gofight.HTTPRequest) {
					response = r
				})

			Convey("When client sends a request to reset the cart", t, func() {
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
					Convey(fmt.Sprintf("Then items, vas-items and item_vas_items tables must be empty"), func() {
						var count int64
						err := TestDB.Table("items").Where("deleted_at IS NULL").Count(&count).Error
						So(err, ShouldBeNil)
						So(count, ShouldEqual, 0)

						err = TestDB.Table("item_vas_items").Where("deleted_at IS NULL").Count(&count).Error
						So(err, ShouldBeNil)
						So(count, ShouldEqual, 0)

						err = TestDB.Table("vas_items").Where("deleted_at IS NULL").Count(&count).Error
						So(err, ShouldBeNil)
						So(count, ShouldEqual, 0)
					})
				}
			})
		})
	}
}
