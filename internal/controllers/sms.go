package controllers

import (
	"net/http"

	"github.com/fastbiztech/hastinapura/internal/services/sms"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func PostSenderID(c *gin.Context) {
	var (
		request = dtos.PostSenderIDRequest{}
		err     error
	)
	// Bind data into the struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = sms.GetService().AddSenderID(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sender ID added successfully"})
}

func GetSenderID(c *gin.Context) {
	var request = dtos.GetSenderIDRequest{}

	// Bind data into the struct
	if err := c.BindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := sms.GetService().GetSenderID(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": response})
}

func DeActivateSenderID(c *gin.Context) {
	var (
		request = dtos.DeleteSenderIDRequest{}
		err     error
	)
	// Bind data into the struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = sms.GetService().DeActivateSenderID(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sender ID deleted successfully"})
}

func ApproveSenderID(c *gin.Context) {
	var (
		request = dtos.ApproveSenderIDRequest{}
		err     error
	)
	// Bind data into the struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = sms.GetService().ApproveSenderID(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sender ID approved successfully"})
}
