package cart

import (
	"checkoutProject/pkg/common/apiresponse"
	"checkoutProject/pkg/common/routing"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CartRouter interface {
	routing.Registerer
}

type cartRouter struct {
	cartController CartController
}

func NewCartRouter(cartController CartController) CartRouter {
	return cartRouter{cartController: cartController}
}

func NewDefaultCartRouter() CartRouter {
	return NewCartRouter(NewDefaultCartController())
}

func (ctr cartRouter) Register(group *gin.RouterGroup) {
	cartGroup := group.Group("")
	cartGroup.GET("", ctr.DisplayCartRoute)
	cartGroup.DELETE("reset", ctr.ResetCartRoute)
}

func (ctr cartRouter) formattedLogger(l logrus.FieldLogger) *logrus.Entry {
	return l.WithField("router", "cart")
}

func (ctr cartRouter) DisplayCartRoute(c *gin.Context) {
	responder, err := ctr.cartController.DisplayCart()
	if err != nil {
		c.JSON(apiresponse.Failed(err))
		return
	}
	c.JSON(apiresponse.OK(responder))
}

func (ctr cartRouter) ResetCartRoute(c *gin.Context) {
	responder, err := ctr.cartController.ResetCart()
	if err != nil {
		c.JSON(apiresponse.Failed(err))
		return
	}
	c.JSON(apiresponse.OK(responder))
}
