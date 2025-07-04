package api

import (
	"github.com/gin-gonic/gin"
	"github.com/krzysu/website-analyzer/internal/database"
	"github.com/krzysu/website-analyzer/internal/worker"
)

func SetupRoutes(router *gin.Engine, db *database.DB, jobQueue chan worker.Job) {
	// Pass the db instance to the handlers
	router.POST("/urls", AddURL(db, jobQueue))
	router.GET("/urls", GetURLs(db))
	router.GET("/urls/:id", GetURL(db))
	router.DELETE("/urls", DeleteURLs(db))
	router.POST("/urls/rerun", RerunURLs(db, jobQueue))
}
