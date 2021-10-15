package nasConvert_test

import (
	"testing"

	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/openapi/models"
	"github.com/smartystreets/goconvey/convey"
)

func TestRequestedNssaiToModels(t *testing.T) {

	testCases := []struct {
		name         string
		requestNssai nasType.RequestedNSSAI
		expected     []models.MappingOfSnssai
	}{
		{
			"Test correctness",
			nasType.RequestedNSSAI{
				Iei: nasMessage.RegistrationRequestRequestedNSSAIType,
				Len: 25,
				Buffer: []uint8{
					0x01, 0x01,
					0x02, 0x01, 0x02,
					0x04, 0x01, 0x01, 0x02, 0x03,
					0x05, 0x01, 0x01, 0x02, 0x03, 0x03,
					0x08, 0x01, 0x11, 0x22, 0x33, 0x04, 0x01, 0x02, 0x03,
				},
			},
			[]models.MappingOfSnssai{
				{
					ServingSnssai: &models.Snssai{
						Sst: 1,
					},
				},
				{
					ServingSnssai: &models.Snssai{
						Sst: 1,
					},
					HomeSnssai: &models.Snssai{
						Sst: 2,
					},
				},
				{
					ServingSnssai: &models.Snssai{
						Sst: 1,
						Sd:  "010203",
					},
				},
				{
					ServingSnssai: &models.Snssai{
						Sst: 1,
						Sd:  "010203",
					},
					HomeSnssai: &models.Snssai{
						Sst: 3,
					},
				},
				{
					ServingSnssai: &models.Snssai{
						Sst: 1,
						Sd:  "112233",
					},
					HomeSnssai: &models.Snssai{
						Sst: 4,
						Sd:  "010203",
					},
				},
			},
		},
		{
			"Test error handling",
			nasType.RequestedNSSAI{
				Iei: nasMessage.RegistrationRequestRequestedNSSAIType,
				Len: 2,
				Buffer: []uint8{
					0x09, 0x01,
				},
			},
			nil,
		},
	}
	convey.Convey("Convert type from nasType.RequestedNSSAI to []models.MappingOfSnssai", t, func() {
		for _, testCase := range testCases {
			modelNssai, err := nasConvert.RequestedNssaiToModels(&testCase.requestNssai)

			convey.Convey(testCase.name, func() {
				convey.So(modelNssai, convey.ShouldResemble, testCase.expected)

				if testCase.name == "Test error handling" {
					convey.So(err, convey.ShouldBeError)
				} else {
					convey.So(err, convey.ShouldBeNil)
				}
			})
		}
	})
}
