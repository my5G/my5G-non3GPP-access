package nasType_test

import (
	"testing"

	"github.com/free5gc/nas/nasType"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPlain5GSNASMessage(t *testing.T) {
	a := nasType.NewPlain5GSNASMessage()
	assert.NotNil(t, a)
}
