package openapi_test

import (
	"free5gc/lib/openapi"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Struct struct {
	A  []string
	S  string
	N  int
	M  map[int]string
	DT *time.Time
}

func TestConvert(t *testing.T) {
	from := map[string]interface{}{
		"A": []string{"Hello", "World"},
		"S": "This is a string",
		"N": 10,
		"M": map[int]interface{}{
			1: "one",
			2: "two",
			3: "three",
		},
		"DT": "2020-05-01T12:04:05+08:00",
	}

	var to Struct

	err := openapi.Convert(from, &to)
	assert.Nil(t, err, "convert failed")

	// check data is correct
	assert.Equal(t, []string{"Hello", "World"}, to.A)
	assert.Equal(t, "This is a string", to.S)
	assert.Equal(t, 10, to.N)
	assert.Equal(t, map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}, to.M)
	expectTime, err := time.Parse(time.RFC3339, "2020-05-01T12:04:05+08:00")
	assert.Equal(t, expectTime, *to.DT)

}
