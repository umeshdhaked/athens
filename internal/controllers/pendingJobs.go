package controllers

import (
	"net/http"

	"github.com/umeshdhaked/athens/internal/services/pendingJobs"
	"github.com/umeshdhaked/athens/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func GetPendingJobs(c *gin.Context) {
	var request = dtos.GetPendingJobsRequest{}

	// Bind the form data into the FormInput struct
	if err := c.BindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := pendingJobs.GetService().GetPendingJobs(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": response})
}
