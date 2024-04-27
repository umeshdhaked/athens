package controllers

import (
	"net/http"

	"github.com/fastbiztech/hastinapura/internal/services/s3Processing"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func GetS3Processing(c *gin.Context) {
	var request = dtos.GetS3ProcessingRequest{}

	// Bind the form data into the FormInput struct
	if err := c.BindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := s3Processing.GetService().GetS3Processing(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": response})
}
