package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func respondPublicError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}

func respondPublicErrorMsg(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func respondInternalError(c *gin.Context, err error, msg string) {
	log.Error().Err(err).Msg(msg)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
}
