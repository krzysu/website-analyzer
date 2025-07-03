package api

import (
	"github.com/gin-gonic/gin"
	"github.com/krzysu/web-crawler/internal/database"
)

func SetupRoutes(router *gin.Engine, db *database.DB) {
	// Pass the db instance to the handlers
	router.POST("/urls", AddURL(db))
	router.GET("/urls", GetURLs(db))
	router.GET("/urls/:id", GetURL(db))
	router.DELETE("/urls", DeleteURLs(db))
	router.POST("/urls/rerun", RerunURLs(db))
}