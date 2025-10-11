package utils

import (
	"os"

	"github.com/umeshdhaked/athens/internal/constants"
)

func GetEnv() string {
	env := os.Getenv(constants.AppEnv)

	// if IsEmpty(env) {
	// 	env = constants.EnvDev
	// }

	return env
}
