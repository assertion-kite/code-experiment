package route

import (
	"github.com/gin-gonic/gin"
)

type RegisterRouteFunc func(*gin.Engine)

func RegisterRoute(rs []RegisterRouteFunc) {
	r := gin.Default()
	for _, routeFunc := range rs {
		routeFunc(r)
	}
	r.Run()
}
