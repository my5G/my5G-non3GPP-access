package rest_api

import (
	"free5gc/lib/openapi/models"
	"free5gc/src/ue/ue_context"
	"free5gc/src/ue/ue_handler/ue_message"
	"github.com/gin-gonic/gin"
	"net/http"
	"free5gc/src/ue/ue_procedures"
	"free5gc/src/ue/logger"
)

type RegistrationData struct {
	AuthenticationMethod models.AuthMethod `json:"authenticationMethod" yaml:"authenticationMethod" bson:"authenticationMethod" mapstructure:"AuthenticationMethod"`
	SupiOrSuci       	string `json:"supiOrSuci" yaml:"supiOrSuci" bson:"supiOrSuci" mapstructure:"SupiOrSuci"`
	K            		string `json:"k" yaml:"k" bson:"k" mapstructure:"K"`
	OpcType				string `json:"opcType" yaml:"opcType" bson:"opcType" mapstructure:"OpcType"`
	Opc          		string `json:"opc" yaml:"opc" bson:"opc" mapstructure:"Opc"`
	PlmnId           	string `json:"plmnid" yaml:"plmnid" bson:"plmnid" mapstructure:"plmnid"`
	ServingNetworkName  string `json:"servingNetworkName" yaml:"servingNetworkName" bson:"servingNetworkName" mapstructure:"ServingNetworkName"`
	Dnn                 string `json:"dnn" yaml:"dnn" bson:"dnn" mapstructure:"dnn"`
	N3IWFIpAddress		string `json:"n3IWFIpAddress" yaml:"n3IWFIpAddress" bson:"n3IWFIpAddress" mapstructure:"N3IWFIpAddress"`
	SNssai				*models.Snssai `json:"SNssai" yaml:"SNssai" bson:"SNssai" mapstructure:"SNssai"`
	IKEBindAddress		string `json:"ikeBindAddress" yaml:"ikeBindAddress" bson:"ikeBindAddress" mapstructure:"IKEBindAddress"`
	PDUSessionID 		uint8 `json:"PDUSessionID" yaml:"PDUSessionID" bson:"PDUSessionID" mapstructure:"PDUSessionID"`
	GREAddress          string `json:"GREAddress" yaml:"GREAddress" bson:"GREAddress" mapstructure:"GREAddress"`
}

func Registration(c *gin.Context) {
	var regData RegistrationData
	err := c.ShouldBindJSON(&regData)
	if err != nil {
		logger.RegistrationLog.Errorln(err)
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	//handlerMsg := udm_message.NewHandlerMessage(udm_message.EventConfirmAuth, req)
	//udm_handler.SendMessage(handlerMsg)
	//rsp := <-handlerMsg.ResponseChan
	//
	//HTTPResponse := rsp.HTTPResponse
	//c.JSON(HTTPResponse.Status, HTTPResponse.Body)

	// Transfer data to UE Context
	ctx := ue_context.UE_Self()
	ctx.AuthenticationMethod = regData.AuthenticationMethod
	ctx.PermanentKeyValue = regData.K
	ctx.SUPIorSUCI = regData.SupiOrSuci
	ctx.OpcType = regData.OpcType
	ctx.Opc = regData.Opc
	//PlmnId
	//ServingNetworkName
	ctx.N3IWFIpAddress = regData.N3IWFIpAddress
	ctx.IKEBindAddress = regData.IKEBindAddress
	ctx.SNssai = regData.SNssai
	ctx.Dnn = regData.Dnn
	ctx.PDUSessionID = regData.PDUSessionID
	ctx.GREIPAddress = regData.GREAddress
	// if

	// TODO: Execute registration procedure as a new Go routine?
	handlerMsg := ue_message.HandlerMessage{
		Event: ue_message.EventRegistrationProcedure,
		Value: nil,
	}
	ue_message.SendMessage(handlerMsg)
	c.String(http.StatusOK, "Registration Procedure triggered.")
}

func Deregistration(c *gin.Context){
	ue_procedures.DeregistrationProcedure(ue_context.UE_Self())
	c.String(http.StatusOK, "Deregistration Procedure executed")
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
	// TODO: Implement here a ping from UE through the 5GC
}

func Info(c *gin.Context) {
	ctx := ue_context.UE_Self()
	c.JSON(http.StatusOK,ctx)
}