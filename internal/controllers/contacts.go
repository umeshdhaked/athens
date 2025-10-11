package controllers

import (
	"net/http"

	"github.com/umeshdhaked/athens/internal/services/contacts"
	"github.com/umeshdhaked/athens/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func GetGroupContacts(c *gin.Context) {
	var request = dtos.GetGroupContactsRequest{}

	// Bind the form data into the FormInput struct
	if err := c.BindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := contacts.GetService().GetContacts(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": response})
}
