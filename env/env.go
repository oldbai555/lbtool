/**
 * @Author: zjj
 * @Date: 2024/4/12
 * @Desc:
**/

package env

import (
	"os"
	"strings"
)

const (
	Dev     = "dev"
	Release = "release"
)

var env string

func GetMode() string {
	env = os.Getenv("LB_TOOL_MODE")
	if env == "" {
		env = Dev
	}
	return env
}

func SetMode(mode string) {
	env = mode
}

func IsRelease() bool {
	return strings.ToLower(GetMode()) == Release
}

func IsDev() bool {
	return strings.ToLower(GetMode()) == Dev
}
