package handler

import (
	"encoding/binary"
	"fmt"
	"free5gc/src/n3iwf/ike/handler"
	"free5gc/src/n3iwf/ike/message"

	//"free5gc/lib/CommonConsumerTestData/UDM/TestGenAuthData"
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasTestpacket"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/nas/security"
	"free5gc/lib/openapi/models"
	//"free5gc/src/ue/ike/handler"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"math/big"
	"time"

	"free5gc/src/ue/context"
	//"errors"
	ike "free5gc/src/ue/ike/message"
	"free5gc/src/ue/logger"
	//"free5gc/src/ue/n3iwf/n3iwf_message"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net"
)

const(
	Ike_InitReq = iota
	Ike_InitRes
)

const(
	Ike_AUTH1 = iota
	Ike_AUTH2
	Ike_AUTH3
	Ike_AUTH4
	Ike_AUTH5
)

var (
	// Log
	ikeLog *logrus.Entry
	//EAP sinaling
	eapIdentifier uint8
	messageID uint32
	sharedKeyExchangeData []byte
	remoteNonce []byte
	attributeType uint16 = ike.AttributeTypeKeyLength
	keyLength uint16 = 256


 	ikePayload ike.IKEPayloadContainer
	eapReq *ike.EAP
	eapExpanded *ike.EAPExpanded
	decodedNAS *nas.Message


    responseSecurityAssociation *ike.SecurityAssociation
	responseTrafficSelectorInitiator *ike.TrafficSelectorInitiator
	responseTrafficSelectorResponder *ike.TrafficSelectorResponder
	responseConfiguration *ike.Configuration
)

func init() {
	ikeLog = logger.IKELog
}

//ueSendInfo *n3iwf_message.UDPSendInfoGroup, message *ike.IKEMessage
func HandleIKESAINIT(udpConn *net.UDPConn, n3iwfAddr, ueAddr *net.UDPAddr , message *ike.IKEMessage  ) {
	ikeLog.Infoln("Handle Initiator IKE_SA_INIT")
	
	// context init variables
	ueSelf := context.UE_Self()

	if  message != nil && message.Flags == ike.ResponseBitCheck {

		//Load Shares Ike Security  Associate
		ikeSecurityAssociation, ok := ueSelf.IKESALoad(message.InitiatorSPI)
		thisUE :=  ikeSecurityAssociation.ThisUE

		if !ok {
			ikeLog.Fatal("Ike security associaiton is not exists for IKE SPI -> " , message.InitiatorSPI)
			return
		}

		for _, ikePayload := range message.Payloads {
			switch ikePayload.Type() {
			case ike.TypeSA:
				ikeLog.Infoln("Get SA payload")
			case ike.TypeKE:
				remotePublicKeyExchangeValue := ikePayload.(*ike.KeyExchange).KeyExchangeData
				var i int = 0
				for {
					if remotePublicKeyExchangeValue[i] != 0 {
						break
					}
				}
				remotePublicKeyExchangeValue = remotePublicKeyExchangeValue[i:]
				remotePublicKeyExchangeValueBig := new(big.Int).SetBytes(remotePublicKeyExchangeValue)
				sharedKeyExchangeData = new(big.Int).Exp(remotePublicKeyExchangeValueBig, thisUE.Secert, thisUE.Factor).Bytes()
			case ike.TypeNiNr:
				remoteNonce = ikePayload.(*ike.Nonce).NonceData
			}
		}

		ikeSecurityAssociation.RemoteSPI = message.ResponderSPI
		ikeSecurityAssociation.ConcatenatedNonce = append(ikeSecurityAssociation.ConcatenatedNonce, remoteNonce...)
		ikeSecurityAssociation.DiffieHellmanSharedKey = sharedKeyExchangeData

		if err := generateKeyForIKESA(ikeSecurityAssociation); err != nil {
			ikeLog.Fatalf("Generate key for IKE SA failed: %+v", err)
		}

	} else {

		//new Security Asssociation
		ikeSecurityAssociation := ueSelf.NewIKESecurityAssociation()
		thisUE := ikeSecurityAssociation.ThisUE

		ikeLog.Infof("Handle Initiator new Ike SA SPI -> %d \n",
			ikeSecurityAssociation.LocalSPI)

		//make IKE Request Msg Sender
		ikeMessage := new(ike.IKEMessage)
		ikeMessage.BuildIKEHeader(ikeSecurityAssociation.LocalSPI, 0, ike.IKE_SA_INIT,
			ike.InitiatorBitCheck, 0)

		// Security Association
		securityAssociation := ikeMessage.Payloads.BuildSecurityAssociation()
		// Proposal 1
		proposal := securityAssociation.Proposals.BuildProposal(1,
			ike.TypeIKE, nil)
		// ENCR
		proposal.EncryptionAlgorithm.BuildTransform(ike.TypeEncryptionAlgorithm, ike.ENCR_AES_CBC,
			&attributeType, &keyLength, nil)
		// INTEG
		proposal.IntegrityAlgorithm.BuildTransform(ike.TypeIntegrityAlgorithm,
			ike.AUTH_HMAC_SHA1_96, nil, nil, nil)
		// PRF
		proposal.PseudorandomFunction.BuildTransform(ike.TypePseudorandomFunction,
			ike.PRF_HMAC_SHA1, nil, nil, nil)
		// DH
		proposal.DiffieHellmanGroup.BuildTransform(ike.TypeDiffieHellmanGroup,
			ike.DH_2048_BIT_MODP, nil, nil, nil)

		// Key exchange data
		generator := new(big.Int).SetUint64(handler.Group14Generator)
		factor, ok := new(big.Int).SetString(handler.Group14PrimeString, 16)
		if !ok {
			ikeLog.Fatal("Generate key exchange datd failed")
		}
		secert := handler.GenerateRandomNumber()
		localPublicKeyExchangeValue := new(big.Int).Exp(generator, secert, factor).Bytes()
		prependZero := make([]byte, len(factor.Bytes())-len(localPublicKeyExchangeValue))
		localPublicKeyExchangeValue = append(prependZero, localPublicKeyExchangeValue...)
		ikeMessage.Payloads.BUildKeyExchange(ike.DH_2048_BIT_MODP, localPublicKeyExchangeValue)

		// Nonce
		localNonce := handler.GenerateRandomNumber().Bytes()
		ikeMessage.Payloads.BuildNonce(localNonce)

		// Store SA  values in virutal Context
		ikeSecurityAssociation.EncryptionAlgorithm = proposal.EncryptionAlgorithm[0]
		ikeSecurityAssociation.IntegrityAlgorithm = proposal.IntegrityAlgorithm[0]
		ikeSecurityAssociation.PseudorandomFunction = proposal.PseudorandomFunction[0]
		ikeSecurityAssociation.DiffieHellmanGroup = proposal.DiffieHellmanGroup[0]
		ikeSecurityAssociation.ConcatenatedNonce = localNonce

		thisUE.Factor = factor
		thisUE.Secert = secert

		//make context IKA Security Association
		ikeLog.Infoln("Send Ike SA Init Request ")
		SendIKEMessageToN3IWF(udpConn, ueAddr, n3iwfAddr, ikeMessage)
	}

}

func HandleIKEAUTH(udpConn *net.UDPConn, n3iwfAddr, ueAddr *net.UDPAddr, ikeMessage *ike.IKEMessage) {

	ikeLog.Infoln("Handle initiator IKE_AUTH ")

	// context init variables
	ueSelf := context.UE_Self()

	/*
		Used to request and reponse value to peer
	*/

	if ikeMessage == nil {
		ikeLog.Error("IKE Message is nil")
		return
	}

	//Load Ike SA
	localSPI := ikeMessage.InitiatorSPI
	ikeSecurityAssociation, ok := ueSelf.IKESALoad(localSPI)
	if !ok {
		ikeLog.Warn("Unrecognized SPI -> ", localSPI)
		return
	}
	//Local Virtual UE
	thisUE := ikeSecurityAssociation.ThisUE

	encryptedPayload, ok := ikeMessage.Payloads[0].(*ike.Encrypted)
	if !ok {
		ikeLog.Fatal("Received payload is not an encrypted payload")
	}

	decryptedIKEPayload, err := decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		ikeLog.Fatalf("Decrypt IKE message failed: %+v", err)
	}

	switch ikeSecurityAssociation.State {
	case Ike_AUTH1:

		// IKE_AUTH Request_3
		ikeMessage.Payloads.Reset()
		ikeMessage.BuildIKEHeader(localSPI, ikeSecurityAssociation.RemoteSPI, ike.IKE_AUTH,
			ike.InitiatorBitCheck, 1)

		// Identification
		ikePayload.BuildIdentificationInitiator(ike.ID_FQDN, []byte("UE"))

		// Security Association
		securityAssociation := ikePayload.BuildSecurityAssociation()
		// Proposal 1
		proposal := securityAssociation.Proposals.BuildProposal(1, ike.TypeESP, []byte{0, 0, 0, 1})
		// ENCR
		proposal.EncryptionAlgorithm.BuildTransform(ike.TypeEncryptionAlgorithm, ike.ENCR_AES_CBC,
			&thisUE.AttributeType, &thisUE.KeyLength, nil)
		// INTEG
		proposal.IntegrityAlgorithm.BuildTransform(ike.TypeIntegrityAlgorithm, ike.AUTH_HMAC_SHA1_96,
			nil, nil, nil)
		// ESN
		proposal.ExtendedSequenceNumbers.BuildTransform(ike.TypeExtendedSequenceNumbers, ike.ESN_NO,
			nil, nil, nil)

		// Traffic Selector
		tsi := ikePayload.BuildTrafficSelectorInitiator()
		tsi.TrafficSelectors.BuildIndividualTrafficSelector(ike.TS_IPV4_ADDR_RANGE, 0, 0, 65535,
			[]byte{0, 0, 0, 0}, []byte{255, 255, 255, 255})
		tsr := ikePayload.BuildTrafficSelectorResponder()
		tsr.TrafficSelectors.BuildIndividualTrafficSelector(ike.TS_IPV4_ADDR_RANGE, 0, 0, 65535,
			[]byte{0, 0, 0, 0}, []byte{255, 255, 255, 255})

		if err := encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage); err != nil {
			ikeLog.Fatalf("Encrypting IKE message failed: %+v", err)
		}

		SendIKEMessageToN3IWF(udpConn, ueAddr, n3iwfAddr, ikeMessage)
		ikeSecurityAssociation.State++

	case Ike_AUTH2:

		// Response 4
		for _, ikePayload := range decryptedIKEPayload {
			switch ikePayload.Type() {
			case ike.TypeIDr:
				ikeLog.Infoln("Get IDr")
			case ike.TypeAUTH:
				ikeLog.Infoln("Get AUTH")
			case ike.TypeCERT:
				ikeLog.Infoln("Get CERT")
			case ike.TypeEAP:
				eapIdentifier = ikePayload.(*ike.EAP).Identifier
				ikeLog.Infoln("Get EAP")
			}
		}
	 	//------------------request 5
		// IKE_AUTH - EAP exchange
		ikeMessage.Payloads.Reset()
		ikeMessage.BuildIKEHeader(localSPI, ikeSecurityAssociation.RemoteSPI, ike.IKE_AUTH,
			ike.InitiatorBitCheck, 2)

		ikePayload.Reset()

		// EAP-5G vendor type data
		eapVendorTypeData := make([]byte, 2)
		eapVendorTypeData[0] = ike.EAP5GType5GNAS

		// AN Parameters
		anParameters := buildEAP5GANParameters()
		anParametersLength := make([]byte, 2)
		binary.BigEndian.PutUint16(anParametersLength, uint16(len(anParameters)))
		eapVendorTypeData = append(eapVendorTypeData, anParametersLength...)
		eapVendorTypeData = append(eapVendorTypeData, anParameters...)

		// NAS
		thisUE.Ran = ike.NewUeRanContext("imsi-2089300007487", 1, security.AlgCiphering128NEA0,
			security.AlgIntegrity128NIA2)

		mobileIdentity5GS := nasType.MobileIdentity5GS{
			Len:    12, // suci
			Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
		}


		ueSecurityCapability := thisUE.Ran.GetUESecurityCapability()
		registrationRequest := nasTestpacket.GetRegistrationRequest(nasMessage.RegistrationType5GSInitialRegistration,
			mobileIdentity5GS, nil, ueSecurityCapability, nil, nil, nil)

		nasLength := make([]byte, 2)
		binary.BigEndian.PutUint16(nasLength, uint16(len(registrationRequest)))
		eapVendorTypeData = append(eapVendorTypeData, nasLength...)
		eapVendorTypeData = append(eapVendorTypeData, registrationRequest...)

		eap := ikePayload.BuildEAP(ike.EAPCodeResponse, eapIdentifier)
		eap.EAPTypeData.BuildEAPExpanded(ike.VendorID3GPP, ike.VendorTypeEAP5G, eapVendorTypeData)

		if err := encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage); err != nil {
			ikeLog.Fatal(err)
		}

		SendIKEMessageToN3IWF(udpConn, ueAddr, n3iwfAddr, ikeMessage)
		ikeSecurityAssociation.State++

	case Ike_AUTH3:

		eapReq, ok = decryptedIKEPayload[0].(*ike.EAP)
		if !ok {
			ikeLog.Fatal("Received packet is not an EAP payload")
		}

		eapExpanded, ok = eapReq.EAPTypeData[0].(*ike.EAPExpanded)
		if !ok {
			ikeLog.Fatal("The EAP data is not an EAP expended.")
		}

		// Decode NAS - Authentication Request
		nasData := eapExpanded.VendorData[4:]
		decodedNAS = new(nas.Message)
		if err := decodedNAS.PlainNasDecode(&nasData); err != nil {
			ikeLog.Fatal(err)
		}

		// Calculate for RES*
		/* assert.NotNil(t, decodedNAS) */

		rand := decodedNAS.AuthenticationRequest.GetRANDValue()
		resStat := thisUE.Ran.DeriveRESstarAndSetKey(thisUE.Ran.AuthenticationSubs, rand[:],
			"5G:mnc093.mcc208.3gppnetwork.org")

		// send NAS Authentication Response
		pdu := nasTestpacket.GetAuthenticationResponse(resStat, "")

		// IKE_AUTH - EAP exchange
		ikeMessage.Payloads.Reset()
		ikeMessage.BuildIKEHeader(localSPI, ikeSecurityAssociation.RemoteSPI, ike.IKE_AUTH, ike.InitiatorBitCheck,
			3)

		ikePayload.Reset()

		// EAP-5G vendor type data
		eapVendorTypeData := make([]byte, 4)
		eapVendorTypeData[0] = ike.EAP5GType5GNAS

		// NAS - Authentication Response
		nasLength := make([]byte, 2)
		binary.BigEndian.PutUint16(nasLength, uint16(len(pdu)))
		eapVendorTypeData = append(eapVendorTypeData, nasLength...)
		eapVendorTypeData = append(eapVendorTypeData, pdu...)

		eap := ikePayload.BuildEAP(ike.EAPCodeResponse, eapReq.Identifier)
		eap.EAPTypeData.BuildEAPExpanded(ike.VendorID3GPP, ike.VendorTypeEAP5G, eapVendorTypeData)

		err := encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage)
		if err != nil {
			ikeLog.Fatal(err)
		}
	case Ike_AUTH4:
		eapReq, ok = decryptedIKEPayload[0].(*ike.EAP)
		if !ok {
			ikeLog.Fatal("Received packet is not an EAP payload")
		}
		eapExpanded, ok = eapReq.EAPTypeData[0].(*ike.EAPExpanded)
		if !ok {
			ikeLog.Fatal("Received packet is not an EAP expended payload")
		}

		//nasData := eapExpanded.VendorData[4:]

		// Send NAS Security Mode Complete Msg
		registrationRequestWith5GMM := nasTestpacket.GetRegistrationRequest(nasMessage.RegistrationType5GSInitialRegistration,
			thisUE.Ran.GetMobileIdentity5GS(), nil, thisUE.Ran.GetUESecurityCapability(),
			thisUE.Ran.Get5GMMCapability(), nil, nil)


		pdu := nasTestpacket.GetSecurityModeComplete(registrationRequestWith5GMM)
		pdu, err := EncodeNasPduWithSecurity(thisUE.Ran,
			pdu,
			nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext,
			true,
			true)

		/* assert.Nil(t, err) */

		// IKE_AUTH - EAP exchange
		ikeMessage.Payloads.Reset()
		ikeMessage.BuildIKEHeader(localSPI, ikeSecurityAssociation.RemoteSPI, ike.IKE_AUTH,
			ike.InitiatorBitCheck, 4)

		ikePayload.Reset()

		// EAP-5G vendor type data
		eapVendorTypeData := make([]byte, 4)
		eapVendorTypeData[0] = ike.EAP5GType5GNAS

		// NAS - Authentication Response
		nasLength := make([]byte, 2)
		binary.BigEndian.PutUint16(nasLength, uint16(len(pdu)))
		eapVendorTypeData = append(eapVendorTypeData, nasLength...)
		eapVendorTypeData = append(eapVendorTypeData, pdu...)

		eap := ikePayload.BuildEAP(ike.EAPCodeResponse, eapReq.Identifier)
		eap.EAPTypeData.BuildEAPExpanded(ike.VendorID3GPP, ike.VendorTypeEAP5G, eapVendorTypeData)

		err = encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage)
		if err != nil {
			ikeLog.Fatal(err)
		}

	case Ike_AUTH5:

		decryptedIKEPayload, err = decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
		if err != nil {
			ikeLog.Fatal(err)
		}
		eapReq, ok = decryptedIKEPayload[0].(*ike.EAP)
		if !ok {
			ikeLog.Fatal("Received packet is not an EAP payload")
		}
		if eapReq.Code != ike.EAPCodeSuccess {
			ikeLog.Fatal("Not Success")
		}

		// IKE_AUTH - Authentication
		ikeMessage.Payloads.Reset()
		ikeMessage.BuildIKEHeader(localSPI, ikeSecurityAssociation.RemoteSPI, ike.IKE_AUTH,
			ike.InitiatorBitCheck, 5)

		ikePayload.Reset()

		// Authentication
		ikePayload.BuildAuthentication(ike.SharedKeyMesageIntegrityCode, []byte{1, 2, 3})

		// Configuration Request
		configurationRequest := ikePayload.BuildConfiguration(ike.CFG_REQUEST)
		configurationRequest.ConfigurationAttribute.BuildConfigurationAttribute(ike.INTERNAL_IP4_ADDRESS, nil)

		err = encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage)
		if err != nil {
			ikeLog.Fatal(err)
		}

	}// End fo switch state machine

}

func HandleCREATECHILDSA(udpConn *net.UDPConn, n3iwfAddr, ueUDPAddr *net.UDPAddr, ikeMessage *ike.IKEMessage) {

	ikeLog.Infoln("Handle CREATE_CHILD_SA")

	if ikeMessage == nil {
		ikeLog.Error("Ike Message is nil")
		return
	}

	// context init variables
	ueSelf := context.UE_Self()
	//Load Ike SA
	localSPI := ikeMessage.InitiatorSPI
	ikeSecurityAssociation, ok := ueSelf.IKESALoad(localSPI)
	if !ok {
		ikeLog.Warn("Unrecognized SPI -> ", localSPI)
		return
	}
	//Local Virtual UE
	thisUE := ikeSecurityAssociation.ThisUE


	//requestIKEMessage := new(ike.IKEMessage)
	n3iwfNASAddr := new(net.TCPAddr)
	ueAddr := new(net.IPNet)

	/* decrypted  Ike packet */
	encryptedPayload, ok := ikeMessage.Payloads[0].(*ike.Encrypted)
	if !ok {
		ikeLog.Fatal("Received packet is not and encrypted payload")
	}
	decryptedIKEPayload, err := decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		ikeLog.Fatal(err)
	}

	for _, ikePayload := range decryptedIKEPayload {
		switch ikePayload.Type() {
		case ike.TypeAUTH:
			ikeLog.Infoln("Get Authentication from N3IWF")
		case ike.TypeSA:
			responseSecurityAssociation = ikePayload.(*ike.SecurityAssociation)
			ikeSecurityAssociation.IKEAuthResponseSA = responseSecurityAssociation
		case ike.TypeTSi:
			responseTrafficSelectorInitiator = ikePayload.(*ike.TrafficSelectorInitiator)
		case ike.TypeTSr:
			responseTrafficSelectorResponder = ikePayload.(*ike.TrafficSelectorResponder)
		case ike.TypeN:
			notification := ikePayload.(*ike.Notification)
			if notification.NotifyMessageType == ike.Vendor3GPPNotifyTypeNAS_IP4_ADDRESS {
				n3iwfNASAddr.IP = net.IPv4(notification.NotificationData[0], notification.NotificationData[1],
					notification.NotificationData[2], notification.NotificationData[3])
			}
			if notification.NotifyMessageType == ike.Vendor3GPPNotifyTypeNAS_TCP_PORT {
				n3iwfNASAddr.Port = int(binary.BigEndian.Uint16(notification.NotificationData))
			}
		case ike.TypeCP:
			responseConfiguration = ikePayload.(*ike.Configuration)
			if responseConfiguration.ConfigurationType == ike.CFG_REPLY {
				for _, configAttr := range responseConfiguration.ConfigurationAttribute {
					if configAttr.Type == ike.INTERNAL_IP4_ADDRESS {
						ueAddr.IP = configAttr.Value
					}
					if configAttr.Type == ike.INTERNAL_IP4_NETMASK {
						ueAddr.Mask = configAttr.Value
					}
				}
			}
		}
	}

	childSecurityAssociationContext, err := thisUE.CreateIKEChildSecurityAssociation(ikeSecurityAssociation.IKEAuthResponseSA)
	if err != nil {
		ikeLog.Fatalf("Create child security association context failed: %+v", err)
		return
	}
	err = ParseIPAddressInformationToChildSecurityAssociation(childSecurityAssociationContext, net.ParseIP(thisUE.N3IWFIpAddress),
		responseTrafficSelectorInitiator.TrafficSelectors[0], responseTrafficSelectorResponder.TrafficSelectors[0])
	if err != nil {
		ikeLog.Fatalf("Parse IP address to child security association failed: %+v", err)
		return
	}
	// Select TCP traffic
	childSecurityAssociationContext.SelectedIPProtocol = unix.IPPROTO_TCP

	if err := GenerateKeyForChildSA(ikeSecurityAssociation, childSecurityAssociationContext); err != nil {
		ikeLog.Fatalf("Generate key for child SA failed: %+v", err)
		return
	}

	// Aplly XFRM rules
	if err = ApplyXFRMRule(true, childSecurityAssociationContext); err != nil {
		ikeLog.Fatalf("Applying XFRM rules failed: %+v", err)
		return
	}

	//make new ipsec and GRe Tunnel
	linkIPSec, err := GetLinkInterfaceByName(thisUE.IPSecIfName)
	if err != nil {
		ikeLog.Fatalf("No link named  %s", thisUE.IPSecIfName)
	}

	linkIPSecAddr := &netlink.Addr{
		IPNet: ueAddr,
	}

	if err := netlink.AddrAdd(linkIPSec, linkIPSecAddr); err != nil {
		ikeLog.Fatalf( fmt.Sprintf("Set %s addr failed: %v",
			thisUE.IPSecIfName , err ))
	}

	localTCPAddr := &net.TCPAddr{
		IP: ueAddr.IP,
	}

	thisUE.TcpConnWithN3IWF, err = net.DialTCP("tcp", localTCPAddr, n3iwfNASAddr)
	if err != nil {
		ikeLog.Fatal(err)
	}

	// send NAS Registration Complete Msg
	pdu := nasTestpacket.GetRegistrationComplete(nil)
	pdu, err = EncodeNasPduWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
		true, false)
	if err != nil {
		ikeLog.Fatal(err)
	}

	SendTCPNASMessageToN3IWF(thisUE.TcpConnWithN3IWF, pdu)

	time.Sleep(500 * time.Millisecond)

	// UE request PDU session setup

	/*
	sNssai := models.Snssai{
		Sst: 1,
		Sd:  "010203",
	}
	*/

	pdu = nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10,
		nasMessage.ULNASTransportRequestTypeInitialRequest,
		thisUE.DNN,
		thisUE.SNssai )

	pdu, err = EncodeNasPduWithSecurity(ue,
		pdu,
		nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
		true,
		false )

	SendTCPNASMessageToN3IWF(thisUE.TcpConnWithN3IWF, pdu)

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		ikeLog.Fatal(err)
	}
	ikeMessage.Payloads.Reset()
	err = ikeMessage.Decode(buffer[:n])
	if err != nil {
		ikeLog.Fatal(err)
	}

	ikeLog.Infof("IKE message exchange type: %d", ikeMessage.ExchangeType)
	ikeLog.Infof("IKE message ID: %d", ikeMessage.MessageID)

	encryptedPayload, ok = ikeMessage.Payloads[0].(*ike.Encrypted)
	if !ok {
		ikeLog.Fatal("Received pakcet is not and encrypted payload")
	}
	decryptedIKEPayload, err = decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		ikeLog.Fatal(err)
	}

	var upIPAddr net.IP
	for _, ikePayload := range decryptedIKEPayload {
		switch ikePayload.Type() {
		case ike.TypeSA:
			responseSecurityAssociation = ikePayload.(*ike.SecurityAssociation)
		case ike.TypeTSi:
			responseTrafficSelectorInitiator = ikePayload.(*ike.TrafficSelectorInitiator)
		case ike.TypeTSr:
			responseTrafficSelectorResponder = ikePayload.(*ike.TrafficSelectorResponder)
		case ike.TypeN:
			notification := ikePayload.(*ike.Notification)
			if notification.NotifyMessageType == ike.Vendor3GPPNotifyType5G_QOS_INFO {
				ikeLog.Infoln("Received Qos Flow settings")
			}
			if notification.NotifyMessageType == ike.Vendor3GPPNotifyTypeUP_IP4_ADDRESS {
				ikeLog.Infof("UP IP Address: %+v\n", notification.NotificationData)
				upIPAddr = notification.NotificationData[:4]
			}
		case ike.TypeNiNr:
			responseNonce := ikePayload.(*ike.Nonce)
			ikeSecurityAssociation.ConcatenatedNonce = responseNonce.NonceData
		}
	}

	// IKE CREATE_CHILD_SA response
	ikeMessage.Payloads.Reset()
	ikeMessage.BuildIKEHeader(ikeMessage.InitiatorSPI, ikeMessage.ResponderSPI, ike.CREATE_CHILD_SA,
		ike.ResponseBitCheck, ikeMessage.MessageID)

	ikePayload.Reset()

	// SA
	ikePayload = append(ikePayload, responseSecurityAssociation)

	// TSi
	ikePayload = append(ikePayload, responseTrafficSelectorInitiator)

	// TSr
	ikePayload = append(ikePayload, responseTrafficSelectorResponder)

	// Nonce
	localNonce := handler.GenerateRandomNumber().Bytes()
	ikeSecurityAssociation.ConcatenatedNonce = append(ikeSecurityAssociation.ConcatenatedNonce, localNonce...)
	ikePayload.BuildNonce(localNonce)

	if err := encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage); err != nil {
		ikeLog.Fatal(err)
	}

	// Send to N3IWF
	ikeMessageData, err := ikeMessage.Encode()
	if err != nil {
		ikeLog.Fatal(err)
	}
	_, err = udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr)
	if err != nil {
		ikeLog.Fatal(err)
	}

	childSecurityAssociationContextUserPlane, err := thisUe.CreateIKEChildSecurityAssociation(responseSecurityAssociation)
	if err != nil {
		ikeLog.Fatalf("Create child security association context failed: %+v", err)
		return
	}
	err = ParseIPAddressInformationToChildSecurityAssociation(childSecurityAssociationContextUserPlane,
		net.ParseIP(thisUE.N3IWFIpAddress), responseTrafficSelectorResponder.TrafficSelectors[0],
		responseTrafficSelectorInitiator.TrafficSelectors[0])
	if err != nil {
		ikeLog.Fatalf("Parse IP address to child security association failed: %+v", err)
		return
	}
	// Select GRE traffic
	childSecurityAssociationContextUserPlane.SelectedIPProtocol = unix.IPPROTO_GRE

	if err = GenerateKeyForChildSA(ikeSecurityAssociation,
		childSecurityAssociationContextUserPlane); err != nil {
		ikeLog.Fatalf("Generate key for child SA failed: %+v", err)
		return
	}

	ikeLog.Logf("State function: encr: %d, auth: %d",
		childSecurityAssociationContextUserPlane.EncryptionAlgorithm,
		childSecurityAssociationContextUserPlane.IntegrityAlgorithm)

	// Aplly XFRM rules
	if err = ApplyXFRMRule(false, childSecurityAssociationContextUserPlane); err != nil {
		ikeLog.Fatalf("Applying XFRM rules failed: %+v", err)
		return
	}

	//Make a new GRE tun interface
	linkGRE, err := SetupNewGRETunnelWithN3IWF(thisUE)

	//Set up address in tun gre interface
	SetGREIpAddr( net.ParseIP(thisUE.GREBindAddress).To4(),
		net.IPv4Mask(255, 255, 255, 0) , linkGRE  )

	//Add Route
	AddGRERoute( linkGRE, net.IPv4zero ,
		net.IPv4Mask(0, 0, 0, 0) )

}