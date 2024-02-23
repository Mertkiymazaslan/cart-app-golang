package routing

import (
	"github.com/gin-gonic/gin"
)

type Registerer interface {
	Register(group *gin.RouterGroup)
}
