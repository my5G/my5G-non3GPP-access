package http

import (
	"free5gc/lib/nas/security"
	"free5gc/lib/openapi/models"
	"free5gc/src/ue/context"
	"free5gc/src/ue/handler/ue_message"
	"github.com/gin-gonic/gin"
	"net/http"
	"free5gc/src/ue/procedures"
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
	N3IWFIpAddress		string `json:"n3IWFIpAddress" yaml:"n3IWFIpAddress" bson:"n3IWFIpAddress" mapstructure:"N3IWFIpAddress"`
	SNssai				*models.Snssai `json:"SNssai" yaml:"SNssai" bson:"SNssai" mapstructure:"SNssai"`
	IKEBindAddress		string `json:"ikeBindAddress" yaml:"ikeBindAddress" bson:"ikeBindAddress" mapstructure:"IKEBindAddress"`
	DNN                 string `json:"DNN" yaml:"DNN" bson:"DNN" mapstructure:"DNN"`
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

	// Transfer data to UE Context
	ctx := context.UE_Self()
	ctx.N3IWFIpAddress = regData.N3IWFIpAddress
	ctx.IKEBindAddress = regData.IKEBindAddress

	/* Make New Virtual UE*/
	newUe := ctx.NewN3iwfUe()
	newUe.CmState = models.CmState_IDLE
	newUe.RmState = models.RmState_DEREGISTERED
	newUe.SUPI = regData.SupiOrSuci
	newUe.DNN = regData.DNN
	newUe.SNssai = regData.SNssai

	/* Make Ike new context Security Association*/
	securityAssociation := ctx.NewIKESecurityAssociation()
	if securityAssociation != nil {
		securityAssociation.ThisUE = newUe
		newUe.LocalSPI = securityAssociation.LocalSPI
	}

	/* Make new RAN UE context */
	newUe.Ran = context.NewUeRanContext(newUe.SUPI, newUe.RanUeNgapId, security.AlgCiphering128NEA0,
		security.AlgIntegrity128NIA2  )

	/* Get Auth Subscription */
	newUe.OpcType = regData.OpcType
	newUe.PermanentKeyValue = regData.K
	newUe.AuthenticationMethod = regData.AuthenticationMethod
	newUe.Opc = regData.Opc
	newUe.Ran.GetAuthSubscription(newUe.PermanentKeyValue, newUe.OpcType, newUe.Opc)

	// if
	// TODO: Execute registration procedure as a new Go routine?
	handlerMsg := ue_message.HandlerMessage{
		Event: ue_message.EventRegistrationProcedure,
		Value: nil,
	}
	//PlmnId
	//ServingNetworkName

	ue_message.SendMessage(handlerMsg)
	c.String(http.StatusOK, "Registration Procedure triggered.")
}

func Deregistration(c *gin.Context){
	procedures.DeregistrationProcedure(context.UE_Self())
	c.String(http.StatusOK, "Deregistration Procedure executed")
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
	// TODO: Implement here a ping from UE through the 5GC
}

func Info(c *gin.Context) {
	ctx := context.UE_Self()
	c.JSON(http.StatusOK,ctx)
}