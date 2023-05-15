package dailylogrus

import (
	"testing"
)

func TestNoInitLog(t *testing.T) {
	Logger.Infoln("hello")
}

func TestInitLogger(t *testing.T) {
	InitLogger("test")
	Logger.Infoln("hello")
}
