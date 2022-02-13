package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func newServer(db storage) (*gin.Engine, error) {
	r := gin.Default()

	r.GET("/latest", func(c *gin.Context) {

		curIP, err := db.latest()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}

		c.JSON(200, curIP)
	})

	r.GET("/history", func(c *gin.Context) {
		history, err := db.history()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}

		c.JSON(200, history)
	})

	return r, nil
}
