package middleware

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	SwaggerInit(r)
}
