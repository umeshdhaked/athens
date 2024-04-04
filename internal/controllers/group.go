package controllers

import (
	"net/http"

	"github.com/fastbiztech/hastinapura/internal/services/group"
	"github.com/fastbiztech/hastinapura/pkg/models/dtos"
	"github.com/gin-gonic/gin"
)

func UploadGroupContacts(c *gin.Context) {
	// Parse multipart form data
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request = dtos.UploadGroupContactsRequest{}
	// Bind the form data into the FormInput struct
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	_, err = group.GetService().UploadGroupToS3(c, file, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded to S3 successfully"})
}
