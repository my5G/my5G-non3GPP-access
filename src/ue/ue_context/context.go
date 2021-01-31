package ue_context

import (
	"crypto/rsa"
	"free5gc/lib/CommonConsumerTestData/UDM/TestGenAuthData"
	"free5gc/lib/openapi/models"
	"free5gc/lib/ngap/ngapType"
	"free5gc/src/ue/logger"
	"free5gc/src/ue/version"
	//"free5gc/src/ue/ue_procedures"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.org/x/net/ipv4"
	"net"
)

var ueContext = UEContext{}
var contextLog *logrus.Entry

func init() {
	// init log
	contextLog = logger.ContextLog

	// init context
	UE_Self().UeId = "1"
	UE_Self().CmState = models.CmState_IDLE
	UE_Self().RmState = models.RmState_DEREGISTERED
	// N3IWF FQDN
	UE_Self().FQDN = ""
	// Network interface mark for xfrm
	UE_Self().Mark = 5
}

type OpcType string

// List of AuthMethod
const (
	OP OpcType  = "5G_AKA"
	OPc OpcType = "EAP_AKA_PRIME"
)

type UEContext struct {
	Version 						 string
	UeId                             string
	HttpIPv4Address					 string
	HttpIPv4Port					 int

	GUAMI *ngapType.GUAMI // connected AMF global identifier TODO: Discutir
	CmState models.CmState // usada no gerenciamento de conexão (interface de sinalização N1 entre UE e AMF) [IDLE|CONNECTED]
	RmState models.RmState // usado no gerenciamento de registro do UE junto ao núcleo [REGISTERED|DEREGISTERED]

	IKESA 		IKESecurityAssociation
	ChildSA     ChildSecurityAssociation

	// Security and Authentication Data
	//UeRanContext *ue_procedures.UeRanContext

	// N3IWF FQDN
	FQDN string

	// security data
	CertificateAuthority []byte
	UECertificate     []byte
	UEPrivateKey      *rsa.PrivateKey

	// Network interface mark for xfrm
	Mark uint32

	// UE local address
	IKEBindAddress      string
	IPSecGatewayAddress string
	GREBindAddress      string
	TCPPort             uint16
	TCPConnection *net.TCPConn

	// Tunnels
	GRETunnel netlink.Link

	// N3IWF N1 interface raw socket
	N1RawSocket *ipv4.RawConn

	// UDP
	UDPSocketAddr *net.UDPAddr

	// subscriber data
	SUPIorSUCI string

	UeAuthenticationCtx *models.UeAuthenticationCtx

	// TODO: #LABORA Include Network Slicing Selection Information??
	SNssai           *models.Snssai
	// TODO: #LABORA Include UE Configuration Information?? (Used in USIM)

	N3IWFIpAddress	string

	// data for registration procedure
	PermanentKeyValue string
	AuthenticationMethod models.AuthMethod
	OpcType string
	Opc string
}

func InitUeContext(){
	// TODO: #LABORA implement context initiation tasks here
	// inicializa com valores default e tudo desregistrado
	ueContext.Version = version.GetVersion()
}

// Reset UE Context
func (context *UEContext) Reset() {
	context.UeId = ""
	// TODO: #LABORA implement context reset here
}

// Create new UE context
func UE_Self() *UEContext {
	return &ueContext
}

func (context *UEContext) SetupUDPSocket(log *logrus.Entry) *net.UDPConn {
	// TODO: #LABORA Add check to verify if socket (udp, port and ip) is already being used
	bindAddr := context.IKEBindAddress + ":500"
	udpAddr, err := net.ResolveUDPAddr("udp", bindAddr)
	if err != nil {
		log.Fatal("Resolve UDP address failed")
	}
	udpListener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("Listen UDP socket failed: %+v", err)
	}
	return udpListener
}

func (context *UEContext) GetAuthSubscription() (authSubs models.AuthenticationSubscription) {
	authSubs.PermanentKey = &models.PermanentKey{
		PermanentKeyValue: context.PermanentKeyValue,
	}
	authSubs.Opc = &models.Opc{
		OpcValue: context.Opc,
	}
	authSubs.Milenage = &models.Milenage{
		Op: &models.Op{
			OpValue: TestGenAuthData.MilenageTestSet19.OP,
		},
	}
	authSubs.AuthenticationManagementField = "8000"

	authSubs.SequenceNumber = TestGenAuthData.MilenageTestSet19.SQN
	authSubs.AuthenticationMethod = context.AuthenticationMethod
	return
}

// temporário
func (context *UEContext) GetAuthSubscription2() (authSubs models.AuthenticationSubscription) {
	authSubs.PermanentKey = &models.PermanentKey{
		PermanentKeyValue: TestGenAuthData.MilenageTestSet19.K,
	}
	authSubs.Opc = &models.Opc{
		OpcValue: TestGenAuthData.MilenageTestSet19.OPC,
	}
	authSubs.Milenage = &models.Milenage{
		Op: &models.Op{
			OpValue: TestGenAuthData.MilenageTestSet19.OP,
		},
	}
	authSubs.AuthenticationManagementField = "8000"

	authSubs.SequenceNumber = TestGenAuthData.MilenageTestSet19.SQN
	authSubs.AuthenticationMethod = models.AuthMethod__5_G_AKA
	return
}