package transport

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) GetById(c *gin.Context) {
	uid := c.Param("id")
	log.Printf("[endpoint][0H] [GetById]: %s", uid)
	data, err := app.orderSvc.Get(uid)
	if err != nil {
		log.Printf("[endpoint][1H] err: %s", err)
		c.JSON(http.StatusNotFound, gin.H{
			"result": "not found",
			"data":   "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"data":   data,
	})
}
