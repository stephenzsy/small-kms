package admin

import (
	"log"

	"github.com/gin-gonic/gin"
)

func respondPublicError(c *gin.Context, status int, err error) {
	log.Printf("Error (%d): %s", status, err.Error())
	c.JSON(status, gin.H{"error": err.Error()})
}
