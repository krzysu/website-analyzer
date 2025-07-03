package api

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine) {
	router.POST("/urls", AddURL)
	router.GET("/urls", GetURLs)
	router.GET("/urls/:id", GetURL)
	router.DELETE("/urls", DeleteURLs)
	router.POST("/urls/rerun", RerunURLs)
}
