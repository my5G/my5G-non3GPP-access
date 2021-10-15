package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMarshalStartTime(t *testing.T) {
	testData := StartTime{
		StartTime: time.Date(1972, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{135, 108, 229, 128}, buf)
}

func TestUnmarshalStartTime(t *testing.T) {
	buf := []byte{135, 108, 229, 128}
	var testData StartTime
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := StartTime{
		StartTime: time.Date(1972, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal(t, expectData, testData)
}
