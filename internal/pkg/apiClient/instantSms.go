package apiClient

import (
	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/internal/pkg/http"
	"github.com/gin-gonic/gin"
)

type InstantSmsApiClient struct {
}

type InstantSmsRespone struct {
}

func (a *InstantSmsApiClient) SendInstantSms(c *gin.Context) (InstantSmsRespone, error) {
	instantSmsApiConfig := config.GetConfig().Api.InstantSms

	_, err := http.NewHTTPClient(instantSmsApiConfig.BaseUrl).
		Method(instantSmsApiConfig.Method).
		Path(instantSmsApiConfig.Path).
		Body(map[string]interface{}{
			// TODO add instant sms body
		}).Headers(map[string]string{
		// TODO add headers
	}).Request(c)

	if err != nil {
		return InstantSmsRespone{}, err
	}

	return InstantSmsRespone{}, nil
}
