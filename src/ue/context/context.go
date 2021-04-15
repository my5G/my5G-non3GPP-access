package context

import (
	"crypto/rand"
	"crypto/rsa"
	"free5gc/lib/CommonConsumerTestData/UDM/TestGenAuthData"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"free5gc/src/ue/logger"
	"free5gc/src/ue/version"
	"git.cs.nctu.edu.tw/calee/sctp"
	gtpv1 "github.com/wmnsk/go-gtp/v1"
	"math"
	"math/big"
	"sync"

	"github.com/free5gc/idgenerator"
	//"free5gc/src/ue/procedures"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/ipv4"
	"net"
)

var ueContext = UEContext{}
var contextLog *logrus.Entry

type OpcType string

// List of AuthMethod
const (
	OP OpcType  = "5G_AKA"
	OPc OpcType = "EAP_AKA_PRIME"
)

type UEContext struct {
	Version  string
	Non3GPPAccessId string
	FQDN string
	HttpIPv4Address	string
	HttpIPv4Port int

	RANUENGAPIDGenerator *idgenerator.IDGenerator

	//Pools of UEs and IPSecs e Security Associate
	UEPool 	sync.Map
	IKESA 	sync.Map
	ChildSA sync.Map

	// UE local address
	IKEBindAddress      string
	UDPSocketAddr *net.UDPAddr
	TCPPort             uint16
	TCPConnection *net.TCPConn

	// N3IWF N1 interface raw socket
	//N1RawSocket *ipv4.RawConn

	UeAuthenticationCtx *models.UeAuthenticationCtx

	// TODO: #LABORA Include UE Configuration Information?? (Used in USIM)
	N3IWFIpAddress	string
}

func init() {
	// init log
	contextLog = logger.ContextLog
	// init context
	UE_Self().Non3GPPAccessId = "1"
	ueContext.Version = version.GetVersion()
	//New IDGenerator
	UE_Self().RANUENGAPIDGenerator = idgenerator.NewGenerator(0, math.MaxInt64)
	// Ue FQND
	UE_Self().FQDN = "UE Non-3GGPP Access"

}

func InitUeContext(){
	// TODO: #LABORA implement context initiation tasks here
	// inicializa com valores default e tudo desregistrado

	//create new Ike Context
	//UE_Self().IKESA = UE_Self().NewIKESecurityAssociation()
}

// Reset UE Context
func (ctx *UEContext) Reset() {
	ctx.UeId = ""
	// TODO: #LABORA implement context reset here
}

// Create new UE context
func UE_Self() *UEContext {
	return &ueContext
}

func (ctx *UEContext) SetupUDPSocket(log *logrus.Entry) *net.UDPConn {
	// TODO: #LABORA Add check to verify if socket (udp, port and ip) is already being used
	bindAddr := ctx.IKEBindAddress + ":500"
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

func (ctx *UEContext) GetAuthSubscription() (authSubs models.AuthenticationSubscription) {
	authSubs.PermanentKey = &models.PermanentKey{
		PermanentKeyValue: ctx.PermanentKeyValue,
	}
	authSubs.Opc = &models.Opc{
		OpcValue: ctx.Opc,
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


func (ctx *UEContext) NewN3iwfUe() *UeN3iwf {

	ranUeNgapId, err := ctx.RANUENGAPIDGenerator.Allocate()
	if err != nil {
		contextLog.Errorf("New UE -> N3IWF failed: %+v", err)
		return nil
	}
	n3iwfUe := new(UeN3iwf)
	n3iwfUe.init(ranUeNgapId)
	ctx.UEPool.Store(ranUeNgapId, n3iwfUe)
	return n3iwfUe
}

func (ctx *UEContext) DeleteN3iwfUe(ranUeNgapId int64) {
	ctx.UEPool.Delete(ranUeNgapId)
}

func (ctx *UEContext) NewIKESecurityAssociation() *IKESecurityAssociation {
	ikeSecurityAssociation := new(IKESecurityAssociation)

	var maxSPI *big.Int = new(big.Int).SetUint64(math.MaxUint64)
	var localSPIuint64 uint64

	for {
		localSPI, err := rand.Int(rand.Reader, maxSPI)
		if err != nil {
			contextLog.Error("[Context] Error occurs when generate new IKE SPI")
			return nil
		}
		localSPIuint64 = localSPI.Uint64()
		if _, duplicate := ctx.IKESA.LoadOrStore(localSPIuint64, ikeSecurityAssociation); !duplicate {
			break
		}
	}

	ikeSecurityAssociation.LocalSPI = localSPIuint64

	return ikeSecurityAssociation
}

func (ctx *UEContext) DeleteIKESecurityAssociation(spi uint64) {
	ctx.IKESA.Delete(spi)
}

func (ctx *UEContext) IKESALoad(spi uint64) (*IKESecurityAssociation, bool) {
	securityAssociation, ok := ctx.IKESA.Load(spi)
	if ok {
		return securityAssociation.(*IKESecurityAssociation), ok
	} else {
		return nil, ok
	}
}