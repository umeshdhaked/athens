package controllers

import (
	"net/http"

	"github.com/umeshdhaked/athens/internal/services/sms"
	"github.com/umeshdhaked/athens/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func PostSms(c *gin.Context) {
	var (
		request = dtos.PostSmsRequest{}
		err     error
	)
	// Bind data into the struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = sms.GetService().SendInstantSms(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sender ID added successfully"})
}
