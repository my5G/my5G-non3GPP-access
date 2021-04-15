package context

import (
	"crypto/rsa"
	"encoding/binary"
	"errors"
	"fmt"
	"free5gc/lib/CommonConsumerTestData/UDM/TestGenAuthData"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	ike_message "free5gc/src/ue/ike/message"
	"free5gc/src/ue/n3iwf/n3iwf_message"
	"github.com/vishvananda/netlink"
	"golang.org/x/net/ipv4"
	"math/big"
	"net"

	gtpv1 "github.com/wmnsk/go-gtp/v1"
)

type UeN3iwf struct {
	SUPI string
	//DNN network
	DNN string
	// data for registration procedure
	AuthenticationMethod models.AuthMethod
	PermanentKeyValue string
	OpcType string
	Opc string

	//Ike SPI
	LocalSPI uint64

	/* UE  RAn identity*/
	RanUeNgapId           int64
	AmfUeNgapId           int64
	PortNumber            int32
	MaskedIMEISV          *ngapType.MaskedIMEISV // TS 38.413 9.3.1.54
	Guti                  string
	RRCEstablishmentCause int16
	GUAMI *ngapType.GUAMI // connected AMF global identifier TODO: Discutir

	//Connection machine state
	CmState models.CmState // usada no gerenciamento de conexão (interface de sinalização N1 entre UE e AMF) [IDLE|CONNECTED]
	RmState models.RmState // usado no gerenciamento de registro do UE junto ao núcleo [REGISTERED|DEREGISTERED]

	// security data
	Factor *big.Int
	Secert *big.Int

	CertificateAuthority []byte
	UECertificate     []byte
	UEPrivateKey      *rsa.PrivateKey

	/* subscriber data generate */
	Ran *UeRanContext

	/* Ip address, Tunnel and interfaces configuration */

	//Interface Name
	IPSecIfName          string
	IKEBindAddress     	 string
	N3IWFIpAddress		 string
	IPSecInnerIP         string
	IPSecGatewayAddress  string

	// xfrm ipsec mark interface
	Mark uint32

	/*  Gre Tunnels */

	//GRE Interfaces
	GREName string
	GREBindAddress   string
	GRENetMask string
	GRETunnel netlink.Link

	/*  Sockets and Ports vars configuration*/
	//Ports NAs port 20000
	TCPNasPort             uint16

	// Sockets

	UDPSocketAddr *net.UDPAddr //For Ike Sender
	TcpConnWithN3IWF *net.TCPConn // For IPsec Nas protocol

	// N3IWF N1 interface raw socket
	N1RawSocket *ipv4.RawConn

	/* Network Slice */
	// TODO: #LABORA Include Network Slicing Selection Information??

	SNssai           *models.Snssai

	/* PDU Session */
	PduSessionList map[int64]*PDUSession // pduSessionId as key

	/* PDU Session Setup Temporary Data */
	TemporaryPDUSessionSetupData *PDUSessionSetupTemporaryData

	/* Security */
	Kn3iwf               []uint8                          // 32 bytes (256 bits), value is from NGAP IE "Security Key"
	SecurityCapabilities *ngapType.UESecurityCapabilities // TS 38.413 9.3.1.86

	/* IKE Security Association */
	N3IWFIKESecurityAssociation   *IKESecurityAssociation
	N3IWFChildSecurityAssociation *ChildSecurityAssociation

	/* NAS IKE Connection */
	UDPSendInfoGroup *n3iwf_message.UDPSendInfoGroup

	/* Others */
	Guami                            *ngapType.GUAMI
	IndexToRfsp                      int64
	Ambr                             *ngapType.UEAggregateMaximumBitRate
	AllowedNssai                     *ngapType.AllowedNSSAI
	RadioCapability                  *ngapType.UERadioCapability                // TODO: This is for RRC, can be deleted
	CoreNetworkAssistanceInformation *ngapType.CoreNetworkAssistanceInformation // TS 38.413 9.3.1.15
	IMSVoiceSupported                int32
}


type PDUSession struct {
	Id                               int64 // PDU Session ID
	Type                             *ngapType.PDUSessionType
	Ambr                             *ngapType.PDUSessionAggregateMaximumBitRate
	Snssai                           ngapType.SNSSAI
	NetworkInstance                  *ngapType.NetworkInstance
	SecurityCipher                   bool
	SecurityIntegrity                bool
	MaximumIntegrityDataRateUplink   *ngapType.MaximumIntegrityProtectedDataRate
	MaximumIntegrityDataRateDownlink *ngapType.MaximumIntegrityProtectedDataRate
	GTPConnection                    *GTPConnectionInfo
	QFIList                          []uint8
	QosFlows                         map[int64]*QosFlow // QosFlowIdentifier as key
}

type PDUSessionSetupTemporaryData struct {
	// Slice of unactivated PDU session
	UnactivatedPDUSession []int64 // PDUSessionID as content
	// NGAPProcedureCode is used to identify which type of
	// response shall be used
	NGAPProcedureCode ngapType.ProcedureCode
	// PDU session setup list response
	SetupListCxtRes  *ngapType.PDUSessionResourceSetupListCxtRes
	FailedListCxtRes *ngapType.PDUSessionResourceFailedToSetupListCxtRes
	SetupListSURes   *ngapType.PDUSessionResourceSetupListSURes
	FailedListSURes  *ngapType.PDUSessionResourceFailedToSetupListSURes
}

type QosFlow struct {
	Identifier int64
	Parameters ngapType.QosFlowLevelQosParameters
}

type GTPConnectionInfo struct {
	UPFIPAddr           string
	UPFUDPAddr          net.Addr
	IncomingTEID        uint32
	OutgoingTEID        uint32
	UserPlaneConnection *gtpv1.UPlaneConn
}

type IKESecurityAssociation struct {
	// SPI
	RemoteSPI uint64
	LocalSPI  uint64

	// Message ID
	MessageID uint32

	// Transforms for IKE SA
	EncryptionAlgorithm    *ike_message.Transform
	PseudorandomFunction   *ike_message.Transform
	IntegrityAlgorithm     *ike_message.Transform
	DiffieHellmanGroup     *ike_message.Transform
	ExpandedSequenceNumber *ike_message.Transform

	// Used for key generating
	ConcatenatedNonce      []byte
	DiffieHellmanSharedKey []byte

	// Keys
	SK_d  []byte // used for child SA key deriving
	SK_ai []byte // used by initiator for integrity checking
	SK_ar []byte // used by responder for integrity checking
	SK_ei []byte // used by initiator for encrypting
	SK_er []byte // used by responder for encrypting
	SK_pi []byte // used by initiator for IKE authentication
	SK_pr []byte // used by responder for IKE authentication

	// State for IKE_AUTH
	State uint8

	// Temporary data stored for the use in later exchange
	InitiatorID              *ike_message.IdentificationInitiator
	InitiatorCertificate     *ike_message.Certificate
	IKEAuthResponseSA        *ike_message.SecurityAssociation
	TrafficSelectorInitiator *ike_message.TrafficSelectorInitiator
	TrafficSelectorResponder *ike_message.TrafficSelectorResponder
	LastEAPIdentifier        uint8

	// Authentication data
	LocalUnsignedAuthentication  []byte
	RemoteUnsignedAuthentication []byte

	// UE context
	ThisUE *UeN3iwf
}

type ChildSecurityAssociation struct {
	// SPI
	SPI uint32

	// IP address
	PeerPublicIPAddr  net.IP
	LocalPublicIPAddr net.IP

	// Traffic selector
	SelectedIPProtocol    uint8
	TrafficSelectorLocal  net.IPNet
	TrafficSelectorRemote net.IPNet

	// Security
	EncryptionAlgorithm               uint16
	InitiatorToResponderEncryptionKey []byte
	ResponderToInitiatorEncryptionKey []byte
	IntegrityAlgorithm                uint16
	InitiatorToResponderIntegrityKey  []byte
	ResponderToInitiatorIntegrityKey  []byte
	ESN                               bool

	// UE context
	ThisUE *UeN3iwf
}

func (ue *UeN3iwf) init() {
	ue.PduSessionList = make(map[int64]*PDUSession)
}

func (ue *UeN3iwf) FindPDUSession(pduSessionID int64) *PDUSession {
	if pduSession, ok := ue.PduSessionList[pduSessionID]; ok {
		return pduSession
	} else {
		return nil
	}
}

func (ue *UeN3iwf) GetAuthSubscription() (authSubs models.AuthenticationSubscription) {
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

func (ue *UeN3iwf) CreatePDUSession(pduSessionID int64, snssai ngapType.SNSSAI) (*PDUSession, error) {
	if _, exists := ue.PduSessionList[pduSessionID]; exists {
		return nil, fmt.Errorf("PDU Session[ID:%d] is already exists", pduSessionID)
	}
	pduSession := &PDUSession{
		Id:       pduSessionID,
		Snssai:   snssai,
		QosFlows: make(map[int64]*QosFlow),
	}
	ue.PduSessionList[pduSessionID] = pduSession
	return pduSession, nil
}


func (ue *UeN3iwf) CreateIKEChildSecurityAssociation(
	chosenSecurityAssociation *ike_message.SecurityAssociation) (*ChildSecurityAssociation, error) {
	childSecurityAssociation := new(ChildSecurityAssociation)

	if chosenSecurityAssociation == nil {
		return nil, errors.New("chosenSecurityAssociation is nil")
	}

	if len(chosenSecurityAssociation.Proposals) == 0 {
		return nil, errors.New("No proposal")
	}

	childSecurityAssociation.SPI = binary.BigEndian.Uint32(chosenSecurityAssociation.Proposals[0].SPI)

	if len(chosenSecurityAssociation.Proposals[0].EncryptionAlgorithm) != 0 {
		childSecurityAssociation.EncryptionAlgorithm = chosenSecurityAssociation.Proposals[0].EncryptionAlgorithm[0].TransformID
	}
	if len(chosenSecurityAssociation.Proposals[0].IntegrityAlgorithm) != 0 {
		childSecurityAssociation.IntegrityAlgorithm = chosenSecurityAssociation.Proposals[0].IntegrityAlgorithm[0].TransformID
	}
	if len(chosenSecurityAssociation.Proposals[0].ExtendedSequenceNumbers) != 0 {
		if chosenSecurityAssociation.Proposals[0].ExtendedSequenceNumbers[0].TransformID == 0 {
			childSecurityAssociation.ESN = false
		} else {
			childSecurityAssociation.ESN = true
		}
	}

	// Link UE context
	childSecurityAssociation.ThisUE = ue
	// Record to N3IWF context
	//n3iwfContext.ChildSA[childSecurityAssociation.SPI] = childSecurityAssociation
	ueContext.ChildSA.Store(childSecurityAssociation.SPI, childSecurityAssociation)

	ue.N3IWFChildSecurityAssociation = childSecurityAssociation

	return childSecurityAssociation, nil
}

func (ue *UeN3iwf) CreateNewSNssaiPDUContext( sst int32, sd string) *ngapType.SNSSAI  {

	ue.SNssai = &models.Snssai{
		Sst: sst,
		Sd:  sd,
	}
	return nil

}
