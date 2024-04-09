package controllers

import (
	"net/http"

	"github.com/fastbiztech/hastinapura/internal/services/sms"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func CreateSmsCampaign(c *gin.Context) {
	var (
		request = dtos.CreateSmsCampaignRequest{}
		err     error
	)
	// Bind data into the struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = sms.GetService().CreateSmsCampaign(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sender ID added successfully"})
}

func GetSmsCampaigns(c *gin.Context) {
	var (
		request = dtos.GetSmsCampaignsRequest{}
		err     error
	)
	// Bind data into the struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = sms.GetService().GetSmsCampaigns(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sender ID added successfully"})
}

func DeActivateSmsCampaigns(c *gin.Context) {
	var (
		request = dtos.DeActivateSmsCampaignRequest{}
		err     error
	)
	// Bind data into the struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = sms.GetService().DeActivateSmsCampaign(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sender ID added successfully"})
}
