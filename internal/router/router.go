package router

import (
    "github.com/gin-gonic/gin"
    "go-api-server/internal/handler"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.GET("/get", handler.GetHandler)
    r.POST("/post", handler.PostHandler)

    return r
}