package controllers

import (
	"net/http"

	"github.com/umeshdhaked/athens/internal/services/sms"
	"github.com/umeshdhaked/athens/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func PostSmsTemplate(c *gin.Context) {
	var (
		request = dtos.PostSmsTemplateRequest{}
		err     error
	)
	// Bind data into the struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = sms.GetService().AddSmsTemplate(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sms Template added successfully"})
}

func GetSmsTemplate(c *gin.Context) {
	var request = dtos.GetSmsTemplateRequest{}

	// Bind data into the struct
	if err := c.BindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := sms.GetService().GetSmsTemplate(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": response})
}

func ApproveSmsTemplate(c *gin.Context) {
	var request = dtos.ApproveSmsTemplateRequest{}

	// Bind data into the struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := sms.GetService().ApproveSmsTemplate(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": response})
}

func UpdateSmsTemplate(c *gin.Context) {
	var request = dtos.UpdateSmsTemplateRequest{}

	// Bind data into the struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := sms.GetService().UpdateSmsTemplate(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": response})
}

func DeActivateSmsTemplate(c *gin.Context) {
	var request = dtos.DeActivateSmsTemplateRequest{}

	// Bind data into the struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := sms.GetService().DeActivateSmsTemplate(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": response})
}
