package middleware

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/xxandjg/ginbase/docs"
	"go.uber.org/zap"
)

func SwaggerInit(r *gin.Engine) {

	zap.L().Info("init Swagger ...")

	// set swagger info
	docs.SwaggerInfo.Title = "ginbase API"
	docs.SwaggerInfo.Description = "https://github.com/xxandjg"
	docs.SwaggerInfo.Version = "v0.0.1"
	docs.SwaggerInfo.Host = "localhost:924"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
