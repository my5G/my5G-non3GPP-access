/*
 * Nudr_DataRepository API OpenAPI file
 *
 * Unified Data Repository Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type ProvisionedDataSets struct {
	AmData      *AccessAndMobilitySubscriptionData  `json:"amData,omitempty" bson:"amData"`
	SmfSelData  *SmfSelectionSubscriptionData       `json:"smfSelData,omitempty" bson:"smfSelData"`
	SmsSubsData *SmsSubscriptionData                `json:"smsSubsData,omitempty" bson:"smsSubsData"`
	SmData      []SessionManagementSubscriptionData `json:"smData,omitempty" bson:"smData"`
	TraceData   *TraceData                          `json:"traceData,omitempty" bson:"traceData"`
	SmsMngData  *SmsManagementSubscriptionData      `json:"smsMngData,omitempty" bson:"smsMngData"`
}
