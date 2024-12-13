package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meta-node-blockchain/meta-node-mns/internal/controller"
)

func InitialRoutes(engine *gin.Engine, controller controller.Controller) {
	r := engine.Group("/api/v1",PreflightHandler())

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	nameRoute := r.Group("/names")
	{
		nameRoute.GET("", controller.GetNames)
		nameRoute.POST("/verify-domain", controller.VerifyDomain)
		nameRoute.GET("/owner/:domain", controller.GetOwnerByName)
	}

}
func PreflightHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
	  c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	  c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	  c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
  
	  if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	  }
	  c.Next()
	}
  }
