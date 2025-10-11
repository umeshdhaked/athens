package controllers

import (
	"net/http"

	"github.com/umeshdhaked/athens/internal/services/sms"
	"github.com/umeshdhaked/athens/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func PostSenderCode(c *gin.Context) {
	var (
		request = dtos.PostSenderCodeRequest{}
		err     error
	)
	// Bind data into the struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = sms.GetService().AddSenderCode(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sender ID added successfully"})
}

func GetSenderCode(c *gin.Context) {
	var request = dtos.GetSenderCodeRequest{}

	// Bind data into the struct
	if err := c.BindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := sms.GetService().GetSenderCode(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": response})
}

func DeActivateSenderCode(c *gin.Context) {
	var (
		request = dtos.DeleteSenderCodeRequest{}
		err     error
	)
	// Bind data into the struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = sms.GetService().DeActivateSenderCode(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sender ID deleted successfully"})
}

func ApproveSenderCode(c *gin.Context) {
	var (
		request = dtos.ApproveSenderCodeRequest{}
		err     error
	)
	// Bind data into the struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = sms.GetService().ApproveSenderCode(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sender ID approved successfully"})
}
