package group

import (
	"fmt"
	"os"
	"time"

	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/pkg/cron"
	"github.com/gin-gonic/gin"
)

type s3ContactsCronExecutor struct {
	ctx *gin.Context
}

func InitialiseS3ContactsCron() {
	// Worker check
	// todo add worker args
	if os.Getenv(constants.WorkerCronArg) != constants.WorkerCronArgS3Contacts {
		return
	}

	if !config.GetConfig().Crons.CronsConfigS3Contacts.Enable {
		return
	}

	newCtx, _ := gin.CreateTestContext(nil)

	job := (&cron.Scheduler{}).NewScheduler()

	job.Initialize(
		time.Duration(config.GetConfig().Crons.CronsConfigS3Contacts.ExecutionTime)*time.Second,
		time.Duration(config.GetConfig().Crons.CronsConfigS3Contacts.StartTime)*time.Second,
		&s3ContactsCronExecutor{ctx: newCtx})
}

func (s *s3ContactsCronExecutor) JobExecutor() {
	// todo Implementation
	fmt.Printf(time.Now().String())
	fmt.Println("i am working boss")
}
