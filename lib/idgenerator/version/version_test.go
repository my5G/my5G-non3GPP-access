package version_test

import (
	"free5gc/lib/idgenerator/version"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, "2020-05-25-01", version.GetVersion())
}
