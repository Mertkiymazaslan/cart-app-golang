package item

import (
	"checkoutProject/pkg/common/apiresponse"
	"checkoutProject/pkg/common/logger"
	"checkoutProject/pkg/common/routing"
	"checkoutProject/pkg/common/validator"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ItemRouter interface {
	routing.Registerer
}

type itemRouter struct {
	itemController ItemController
}

func NewItemRouter(itemController ItemController) ItemRouter {
	return itemRouter{itemController: itemController}
}

func NewDefaultItemRouter() ItemRouter {
	return NewItemRouter(NewDefaultItemController())
}

func (itr itemRouter) Register(group *gin.RouterGroup) {
	itemGroup := group.Group("items")
	itemGroup.POST("", itr.AddItemRoute)
	itemGroup.DELETE(":item_id", itr.RemoveItemRoute)
}

func (itr itemRouter) formattedLogger(l logrus.FieldLogger) *logrus.Entry {
	return l.WithField("router", "item")
}

func (itr itemRouter) AddItemRoute(c *gin.Context) {
	log := itr.formattedLogger(logger.GetInstance()).WithField("location", "AddItemRoute")

	var params AddItemParams

	if err := c.ShouldBindJSON(&params); err != nil {
		readableErr := validator.GetValidatorMessages(err)
		log.WithError(readableErr).Error("Could not bind parameters")
		c.JSON(apiresponse.Failed(readableErr))
		return
	}

	responder, err := itr.itemController.AddItem(params)
	if err != nil {
		c.JSON(apiresponse.Failed(err))
		return
	}

	c.JSON(apiresponse.Created(responder))
}

func (itr itemRouter) RemoveItemRoute(c *gin.Context) {
	log := itr.formattedLogger(logger.GetInstance()).WithField("location", "RemoveItemRoute")

	var params RemoveItemParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.WithError(err).Error("Could not bind parameters")
		c.JSON(apiresponse.Failed(err))
		return
	}

	responder, err := itr.itemController.RemoveItem(params)
	if err != nil {
		c.JSON(apiresponse.Failed(err))
		return
	}

	c.JSON(apiresponse.OK(responder))
}

type VasItemRouter interface {
	routing.Registerer
}

type vasItemRouter struct {
	vasItemController VasItemController
}

func NewVasItemRouter(vasItemController VasItemController) VasItemRouter {
	return vasItemRouter{vasItemController: vasItemController}
}

func NewDefaultVasItemRouter() VasItemRouter {
	return NewVasItemRouter(NewDefaultVasItemController())
}

func (vitr vasItemRouter) Register(group *gin.RouterGroup) {
	vasItemGroup := group.Group("items/:item_id/vas-items")
	vasItemGroup.POST("", vitr.AddVasItemRoute)
}

func (vitr vasItemRouter) formattedLogger(l logrus.FieldLogger) *logrus.Entry {
	return l.WithField("router", "vas-item")
}

func (vitr vasItemRouter) AddVasItemRoute(c *gin.Context) {
	log := vitr.formattedLogger(logger.GetInstance()).WithField("location", "AddVasItemRoute")

	var params AddVasItemParams

	if err := c.ShouldBindUri(&params.ItemUriParams); err != nil {
		readableErr := validator.GetValidatorMessages(err)
		log.WithError(readableErr).Error("Could not bind parameters")
		c.JSON(apiresponse.Failed(readableErr))
		return
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		readableErr := validator.GetValidatorMessages(err)
		log.WithError(readableErr).Error("Could not bind parameters")
		c.JSON(apiresponse.Failed(readableErr))
		return
	}

	responder, err := vitr.vasItemController.AddVasItem(params)
	if err != nil {
		c.JSON(apiresponse.Failed(err))
		return
	}

	c.JSON(apiresponse.Created(responder))
}
