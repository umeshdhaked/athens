package di

import "github.com/FastBizTech/hastinapura/api/services"

// var dynamoConnection pkg.DynnamoConnection
var regService *services.RegistrationService

func InitialiseServices() {
	// dynamoConnection = pkg.NewDynnamoConnection()
	regService = services.NewRegistrationService()
}

func GetRegistrationService() *services.RegistrationService {
	return regService
}
