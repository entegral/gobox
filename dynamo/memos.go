package dynamo

import (
	"os"
)

var isTesting bool
var checkedTesting bool

func checkTesting() bool {
	if checkedTesting {
		return isTesting
	}
	isTesting = os.Getenv("GOBOX_TESTING") != ""
	checkedTesting = true
	return isTesting
}
