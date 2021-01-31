package path_util

import (
	"free5gc/lib/path_util/logger"
	"testing"
)

func TestFree5gcPath(t *testing.T) {
	logger.PathLog.Infoln(Gofree5gcPath("free5gc/abcdef/abcdef.pem"))
}
