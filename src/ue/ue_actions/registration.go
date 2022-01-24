package ue_actions

import (
	"encoding/binary"
	"errors"
	"fmt"
	"free5gc/src/ue/logger"
	"free5gc/src/ue/ue_context"
	"free5gc/src/ue/ue_nas"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"hash"
	"math/big"
	"net"
	"time"

	"github.com/free5gc/CommonConsumerTestData/UDM/TestGenAuthData"
	"github.com/free5gc/n3iwf/context"
	"github.com/free5gc/n3iwf/ike/handler"
	"github.com/free5gc/n3iwf/ike/message"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/openapi/models"
)

var pingLog *logrus.Entry

func init() {
	pingLog = logger.RunLog
}

func createIKEChildSecurityAssociation(chosenSecurityAssociation *message.SecurityAssociation) (*context.ChildSecurityAssociation, error) {
	childSecurityAssociation := new(context.ChildSecurityAssociation)

	if chosenSecurityAssociation == nil {
		return nil, errors.New("chosenSecurityAssociation is nil")
	}

	if len(chosenSecurityAssociation.Proposals) == 0 {
		return nil, errors.New("no proposal")
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

	return childSecurityAssociation, nil
}

func getAuthSubscription() (authSubs models.AuthenticationSubscription) {
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

func getAuthSubscription2() (authSubs models.AuthenticationSubscription) {
	authSubs.PermanentKey = &models.PermanentKey{
		PermanentKeyValue: MilenageTestSet18.K,
	}
	authSubs.Opc = &models.Opc{
		OpcValue: MilenageTestSet18.OPC,
	}
	authSubs.Milenage = &models.Milenage{
		Op: &models.Op{
			OpValue: MilenageTestSet18.OP,
		},
	}
	authSubs.AuthenticationManagementField = "8000"

	authSubs.SequenceNumber = MilenageTestSet18.SQN
	authSubs.AuthenticationMethod = models.AuthMethod__5_G_AKA
	return
}

func setupUDPSocket(ctx *ue_context.UEContext, log *logrus.Entry) *net.UDPConn {
	bindAddr := fmt.Sprintf("%s:500", ctx.IKEBindAddress)
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

func concatenateNonceAndSPI(nonce []byte, SpiInitiator uint64, SpiResponder uint64) []byte {
	spi := make([]byte, 8)

	binary.BigEndian.PutUint64(spi, SpiInitiator)
	newSlice := append(nonce, spi...)
	binary.BigEndian.PutUint64(spi, SpiResponder)
	newSlice = append(newSlice, spi...)

	return newSlice
}

func generateKeyForIKESA(ikeSecurityAssociation *context.IKESecurityAssociation) error {
	// Transforms
	transformPseudorandomFunction := ikeSecurityAssociation.PseudorandomFunction

	// Get key length of SK_d, SK_ai, SK_ar, SK_ei, SK_er, SK_pi, SK_pr
	var lengthSkD, lengthSkAi, lengthSkAr, lengthSkEi, lengthSkEr, lengthSkPi, lengthSkPr, totalKeyLength int
	var ok bool

	lengthSkD = 20
	lengthSkAi = 20
	lengthSkAr = lengthSkAi
	lengthSkEi = 32
	lengthSkEr = lengthSkEi
	lengthSkPi, lengthSkPr = lengthSkD, lengthSkD
	totalKeyLength = lengthSkD + lengthSkAi + lengthSkAr + lengthSkEi + lengthSkEr + lengthSkPi + lengthSkPr

	// Generate IKE SA key as defined in RFC7296 Section 1.3 and Section 1.4
	var pseudorandomFunction hash.Hash

	if pseudorandomFunction, ok = handler.NewPseudorandomFunction(ikeSecurityAssociation.ConcatenatedNonce, transformPseudorandomFunction.TransformID); !ok {
		return errors.New("new pseudorandom function failed")
	}

	if _, err := pseudorandomFunction.Write(ikeSecurityAssociation.DiffieHellmanSharedKey); err != nil {
		return errors.New("pseudorandom function write failed")
	}

	SKEYSEED := pseudorandomFunction.Sum(nil)

	seed := concatenateNonceAndSPI(ikeSecurityAssociation.ConcatenatedNonce, ikeSecurityAssociation.LocalSPI, ikeSecurityAssociation.RemoteSPI)

	var keyStream, generatedKeyBlock []byte
	var index byte
	for index = 1; len(keyStream) < totalKeyLength; index++ {
		if pseudorandomFunction, ok = handler.NewPseudorandomFunction(SKEYSEED, transformPseudorandomFunction.TransformID); !ok {
			return errors.New("new pseudorandom function failed")
		}
		if _, err := pseudorandomFunction.Write(append(append(generatedKeyBlock, seed...), index)); err != nil {
			return errors.New("pseudorandom function write failed")
		}
		generatedKeyBlock = pseudorandomFunction.Sum(nil)
		keyStream = append(keyStream, generatedKeyBlock...)
	}

	// Assign keys into context
	ikeSecurityAssociation.SK_d = keyStream[:lengthSkD]
	keyStream = keyStream[lengthSkD:]
	ikeSecurityAssociation.SK_ai = keyStream[:lengthSkAi]
	keyStream = keyStream[lengthSkAi:]
	ikeSecurityAssociation.SK_ar = keyStream[:lengthSkAr]
	keyStream = keyStream[lengthSkAr:]
	ikeSecurityAssociation.SK_ei = keyStream[:lengthSkEi]
	keyStream = keyStream[lengthSkEi:]
	ikeSecurityAssociation.SK_er = keyStream[:lengthSkEr]
	keyStream = keyStream[lengthSkEr:]
	ikeSecurityAssociation.SK_pi = keyStream[:lengthSkPi]
	keyStream = keyStream[lengthSkPi:]
	ikeSecurityAssociation.SK_pr = keyStream[:lengthSkPr]
	keyStream = keyStream[lengthSkPr:]

	return nil
}

func generateKeyForChildSA(ikeSecurityAssociation *context.IKESecurityAssociation, childSecurityAssociation *context.ChildSecurityAssociation) error {
	// Transforms
	transformPseudorandomFunction := ikeSecurityAssociation.PseudorandomFunction
	var transformIntegrityAlgorithmForIPSec *message.Transform
	if len(ikeSecurityAssociation.IKEAuthResponseSA.Proposals[0].IntegrityAlgorithm) != 0 {
		transformIntegrityAlgorithmForIPSec = ikeSecurityAssociation.IKEAuthResponseSA.Proposals[0].IntegrityAlgorithm[0]
	}

	// Get key length for encryption and integrity key for IPSec
	var lengthEncryptionKeyIPSec, lengthIntegrityKeyIPSec, totalKeyLength int
	var ok bool

	lengthEncryptionKeyIPSec = 32
	if transformIntegrityAlgorithmForIPSec != nil {
		lengthIntegrityKeyIPSec = 20
	}
	totalKeyLength = lengthEncryptionKeyIPSec + lengthIntegrityKeyIPSec
	totalKeyLength = totalKeyLength * 2

	// Generate key for child security association as specified in RFC 7296 section 2.17
	seed := ikeSecurityAssociation.ConcatenatedNonce
	var pseudorandomFunction hash.Hash

	var keyStream, generatedKeyBlock []byte
	var index byte
	for index = 1; len(keyStream) < totalKeyLength; index++ {
		if pseudorandomFunction, ok = handler.NewPseudorandomFunction(ikeSecurityAssociation.SK_d, transformPseudorandomFunction.TransformID); !ok {
			return errors.New("new pseudorandom function failed")
		}
		if _, err := pseudorandomFunction.Write(append(append(generatedKeyBlock, seed...), index)); err != nil {
			return errors.New("pseudorandom function write failed")
		}
		generatedKeyBlock = pseudorandomFunction.Sum(nil)
		keyStream = append(keyStream, generatedKeyBlock...)
	}

	childSecurityAssociation.InitiatorToResponderEncryptionKey = append(childSecurityAssociation.InitiatorToResponderEncryptionKey, keyStream[:lengthEncryptionKeyIPSec]...)
	keyStream = keyStream[lengthEncryptionKeyIPSec:]
	childSecurityAssociation.InitiatorToResponderIntegrityKey = append(childSecurityAssociation.InitiatorToResponderIntegrityKey, keyStream[:lengthIntegrityKeyIPSec]...)
	keyStream = keyStream[lengthIntegrityKeyIPSec:]
	childSecurityAssociation.ResponderToInitiatorEncryptionKey = append(childSecurityAssociation.ResponderToInitiatorEncryptionKey, keyStream[:lengthEncryptionKeyIPSec]...)
	keyStream = keyStream[lengthEncryptionKeyIPSec:]
	childSecurityAssociation.ResponderToInitiatorIntegrityKey = append(childSecurityAssociation.ResponderToInitiatorIntegrityKey, keyStream[:lengthIntegrityKeyIPSec]...)

	return nil

}

func decryptProcedure(ikeSecurityAssociation *context.IKESecurityAssociation, ikeMessage *message.IKEMessage, encryptedPayload *message.Encrypted) (message.IKEPayloadContainer, error) {
	// Load needed information
	transformIntegrityAlgorithm := ikeSecurityAssociation.IntegrityAlgorithm
	transformEncryptionAlgorithm := ikeSecurityAssociation.EncryptionAlgorithm
	checksumLength := 12 // HMAC_SHA1_96

	// Checksum
	checksum := encryptedPayload.EncryptedData[len(encryptedPayload.EncryptedData)-checksumLength:]

	ikeMessageData, err := ikeMessage.Encode()
	if err != nil {
		return nil, errors.New("encoding IKE message failed")
	}

	ok, err := handler.VerifyIKEChecksum(ikeSecurityAssociation.SK_ar, ikeMessageData[:len(ikeMessageData)-checksumLength], checksum, transformIntegrityAlgorithm.TransformID)
	if err != nil {
		return nil, errors.New("error verify checksum")
	}
	if !ok {
		return nil, errors.New("checksum failed, drop")
	}

	// Decrypt
	encryptedData := encryptedPayload.EncryptedData[:len(encryptedPayload.EncryptedData)-checksumLength]
	plainText, err := handler.DecryptMessage(ikeSecurityAssociation.SK_er, encryptedData, transformEncryptionAlgorithm.TransformID)
	if err != nil {
		return nil, errors.New("error decrypting message")
	}

	var decryptedIKEPayload message.IKEPayloadContainer
	err = decryptedIKEPayload.Decode(encryptedPayload.NextPayload, plainText)
	if err != nil {
		return nil, errors.New("decoding decrypted payload failed")
	}

	return decryptedIKEPayload, nil

}

func encryptProcedure(ikeSecurityAssociation *context.IKESecurityAssociation, ikePayload message.IKEPayloadContainer, responseIKEMessage *message.IKEMessage) error {
	// Load needed information
	transformIntegrityAlgorithm := ikeSecurityAssociation.IntegrityAlgorithm
	transformEncryptionAlgorithm := ikeSecurityAssociation.EncryptionAlgorithm
	checksumLength := 12 // HMAC_SHA1_96

	// Encrypting
	notificationPayloadData, err := ikePayload.Encode()
	if err != nil {
		return errors.New("encoding IKE payload failed")
	}

	encryptedData, err := handler.EncryptMessage(ikeSecurityAssociation.SK_ei, notificationPayloadData, transformEncryptionAlgorithm.TransformID)
	if err != nil {
		return errors.New("error encrypting message")
	}

	encryptedData = append(encryptedData, make([]byte, checksumLength)...)
	sk := responseIKEMessage.Payloads.BuildEncrypted(ikePayload[0].Type(), encryptedData)

	// Calculate checksum
	responseIKEMessageData, err := responseIKEMessage.Encode()
	if err != nil {
		return errors.New("encoding IKE message error")
	}
	checksumOfMessage, err := handler.CalculateChecksum(ikeSecurityAssociation.SK_ai, responseIKEMessageData[:len(responseIKEMessageData)-checksumLength], transformIntegrityAlgorithm.TransformID)
	if err != nil {
		return errors.New("error calculating checksum")
	}
	checksumField := sk.EncryptedData[len(sk.EncryptedData)-checksumLength:]
	copy(checksumField, checksumOfMessage)

	return nil

}

// [TS 24502] 9.3.2.2.2 EAP-Response/5G-NAS message
// Define EAP-Response/5G-NAS message and AN-Parameters Format.

// [TS 24501] 8.2.6.1.1  REGISTRATION REQUEST message content
// For dealing with EAP-5G start, return EAP-5G response including
// "AN-Parameters and NASPDU of Registration Request"

func buildEAP5GANParameters() []byte {
	var anParameters []byte

	// [TS 24.502] 9.3.2.2.2.3
	// AN-parameter value field in GUAMI, PLMN ID and NSSAI is coded as value part
	// Therefore, IEI of AN-parameter is not needed to be included.

	// anParameter = AN-parameter Type | AN-parameter Length | Value part of IE

	// Build GUAMI
	anParameter := make([]byte, 2)
	guami := make([]byte, 6)
	guami[0] = 0x02
	guami[1] = 0xf8
	guami[2] = 0x39
	guami[3] = 0xca
	guami[4] = 0xfe
	guami[5] = 0x0
	anParameter[0] = message.ANParametersTypeGUAMI
	anParameter[1] = byte(len(guami))
	anParameter = append(anParameter, guami...)

	anParameters = append(anParameters, anParameter...)

	// Build Establishment Cause
	anParameter = make([]byte, 2)
	establishmentCause := make([]byte, 1)
	establishmentCause[0] = message.EstablishmentCauseMO_Signalling
	anParameter[0] = message.ANParametersTypeEstablishmentCause
	anParameter[1] = byte(len(establishmentCause))
	anParameter = append(anParameter, establishmentCause...)

	anParameters = append(anParameters, anParameter...)

	// Build PLMN ID
	anParameter = make([]byte, 2)
	plmnID := make([]byte, 3)
	plmnID[0] = 0x02
	plmnID[1] = 0xf8
	plmnID[2] = 0x39
	anParameter[0] = message.ANParametersTypeSelectedPLMNID
	anParameter[1] = byte(len(plmnID))
	anParameter = append(anParameter, plmnID...)

	anParameters = append(anParameters, anParameter...)

	// Build NSSAI
	anParameter = make([]byte, 2)
	var nssai []byte
	// s-nssai = s-nssai length(1 byte) | SST(1 byte) | SD(3 bytes)
	snssai := make([]byte, 5)
	snssai[0] = 4
	snssai[1] = 1
	snssai[2] = 0x01
	snssai[3] = 0x02
	snssai[4] = 0x03
	nssai = append(nssai, snssai...)
	snssai = make([]byte, 5)
	snssai[0] = 4
	snssai[1] = 1
	snssai[2] = 0x11
	snssai[3] = 0x22
	snssai[4] = 0x33
	nssai = append(nssai, snssai...)
	anParameter[0] = message.ANParametersTypeRequestedNSSAI
	anParameter[1] = byte(len(nssai))
	anParameter = append(anParameter, nssai...)

	anParameters = append(anParameters, anParameter...)

	return anParameters
}

func parseIPAddressInformationToChildSecurityAssociation(
	childSecurityAssociation *context.ChildSecurityAssociation,
	n3iwfPublicIPAddr net.IP,
	ikeBindAddr string,
	trafficSelectorLocal *message.IndividualTrafficSelector,
	trafficSelectorRemote *message.IndividualTrafficSelector) error {

	if childSecurityAssociation == nil {
		return errors.New("childSecurityAssociation is nil")
	}

	childSecurityAssociation.PeerPublicIPAddr = n3iwfPublicIPAddr
	childSecurityAssociation.LocalPublicIPAddr = net.ParseIP(ikeBindAddr)

	childSecurityAssociation.TrafficSelectorLocal = net.IPNet{
		IP:   trafficSelectorLocal.StartAddress,
		Mask: []byte{255, 255, 255, 255},
	}

	childSecurityAssociation.TrafficSelectorRemote = net.IPNet{
		IP:   trafficSelectorRemote.StartAddress,
		Mask: []byte{255, 255, 255, 255},
	}

	return nil
}

func applyXFRMRule(ueIsInitiator bool, childSecurityAssociation *context.ChildSecurityAssociation) error {
	// Build XFRM information data structure for incoming traffic.

	// Mark
	mark := &netlink.XfrmMark{
		Value: 5,
	}

	// Direction: N3IWF -> UE
	// State
	var xfrmEncryptionAlgorithm, xfrmIntegrityAlgorithm *netlink.XfrmStateAlgo
	if ueIsInitiator {
		xfrmEncryptionAlgorithm = &netlink.XfrmStateAlgo{
			Name: handler.XFRMEncryptionAlgorithmType(childSecurityAssociation.EncryptionAlgorithm).String(),
			Key:  childSecurityAssociation.ResponderToInitiatorEncryptionKey,
		}
		if childSecurityAssociation.IntegrityAlgorithm != 0 {
			xfrmIntegrityAlgorithm = &netlink.XfrmStateAlgo{
				Name: handler.XFRMIntegrityAlgorithmType(childSecurityAssociation.IntegrityAlgorithm).String(),
				Key:  childSecurityAssociation.ResponderToInitiatorIntegrityKey,
			}
		}
	} else {
		xfrmEncryptionAlgorithm = &netlink.XfrmStateAlgo{
			Name: handler.XFRMEncryptionAlgorithmType(childSecurityAssociation.EncryptionAlgorithm).String(),
			Key:  childSecurityAssociation.InitiatorToResponderEncryptionKey,
		}
		if childSecurityAssociation.IntegrityAlgorithm != 0 {
			xfrmIntegrityAlgorithm = &netlink.XfrmStateAlgo{
				Name: handler.XFRMIntegrityAlgorithmType(childSecurityAssociation.IntegrityAlgorithm).String(),
				Key:  childSecurityAssociation.InitiatorToResponderIntegrityKey,
			}
		}
	}

	xfrmState := new(netlink.XfrmState)

	xfrmState.Src = childSecurityAssociation.PeerPublicIPAddr
	xfrmState.Dst = childSecurityAssociation.LocalPublicIPAddr
	xfrmState.Proto = netlink.XFRM_PROTO_ESP
	xfrmState.Mode = netlink.XFRM_MODE_TUNNEL
	xfrmState.Spi = int(childSecurityAssociation.SPI)
	xfrmState.Mark = mark
	xfrmState.Auth = xfrmIntegrityAlgorithm
	xfrmState.Crypt = xfrmEncryptionAlgorithm
	xfrmState.ESN = childSecurityAssociation.ESN

	// Commit xfrm state to netlink
	var err error
	if err = netlink.XfrmStateAdd(xfrmState); err != nil {
		return fmt.Errorf("set XFRM state rule failed: %+v", err)
	}

	// Policy
	xfrmPolicyTemplate := netlink.XfrmPolicyTmpl{
		Src:   xfrmState.Src,
		Dst:   xfrmState.Dst,
		Proto: xfrmState.Proto,
		Mode:  xfrmState.Mode,
		Spi:   xfrmState.Spi,
	}

	xfrmPolicy := new(netlink.XfrmPolicy)

	if childSecurityAssociation.SelectedIPProtocol == 0 {
		return errors.New("protocol == 0")
	}

	xfrmPolicy.Src = &childSecurityAssociation.TrafficSelectorRemote
	xfrmPolicy.Dst = &childSecurityAssociation.TrafficSelectorLocal
	xfrmPolicy.Proto = netlink.Proto(childSecurityAssociation.SelectedIPProtocol)
	xfrmPolicy.Dir = netlink.XFRM_DIR_IN
	xfrmPolicy.Mark = mark
	xfrmPolicy.Tmpls = []netlink.XfrmPolicyTmpl{
		xfrmPolicyTemplate,
	}

	// Commit xfrm policy to netlink
	if err = netlink.XfrmPolicyAdd(xfrmPolicy); err != nil {
		return fmt.Errorf("set XFRM policy rule failed: %+v", err)
	}

	// Direction: UE -> N3IWF
	// State
	if ueIsInitiator {
		xfrmEncryptionAlgorithm.Key = childSecurityAssociation.InitiatorToResponderEncryptionKey
		if childSecurityAssociation.IntegrityAlgorithm != 0 {
			xfrmIntegrityAlgorithm.Key = childSecurityAssociation.InitiatorToResponderIntegrityKey
		}
	} else {
		xfrmEncryptionAlgorithm.Key = childSecurityAssociation.ResponderToInitiatorEncryptionKey
		if childSecurityAssociation.IntegrityAlgorithm != 0 {
			xfrmIntegrityAlgorithm.Key = childSecurityAssociation.ResponderToInitiatorIntegrityKey
		}
	}

	xfrmState.Src, xfrmState.Dst = xfrmState.Dst, xfrmState.Src

	// Commit xfrm state to netlink
	if err = netlink.XfrmStateAdd(xfrmState); err != nil {
		return fmt.Errorf("set XFRM state rule failed: %+v", err)
	}

	// Policy
	xfrmPolicyTemplate.Src, xfrmPolicyTemplate.Dst = xfrmPolicyTemplate.Dst, xfrmPolicyTemplate.Src

	xfrmPolicy.Src, xfrmPolicy.Dst = xfrmPolicy.Dst, xfrmPolicy.Src
	xfrmPolicy.Dir = netlink.XFRM_DIR_OUT
	xfrmPolicy.Tmpls = []netlink.XfrmPolicyTmpl{
		xfrmPolicyTemplate,
	}

	// Commit xfrm policy to netlink
	if err = netlink.XfrmPolicyAdd(xfrmPolicy); err != nil {
		return fmt.Errorf("set XFRM policy rule failed: %+v", err)
	}

	return nil
}

func InitialRegistrationProcedure(ueContext *ue_context.UEContext) {
	// New UE
	ue := NewRanUeContext(fmt.Sprintf("imsi-%s", ueContext.SUPIorSUCI), 1, security.AlgCiphering128NEA0, security.AlgIntegrity128NIA2,
		models.AccessType_NON_3_GPP_ACCESS)

	var UE string
	var mobileIdentity5GS nasType.MobileIdentity5GS

	if fmt.Sprintf("imsi-%s", ueContext.SUPIorSUCI) == "imsi-2089300007486" {
		mobileIdentity5GS = nasType.MobileIdentity5GS{
			Len:    12, // suci
			Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x47, 0x68},
		}
		ue.AuthenticationSubs = getAuthSubscription2()
		ue.AmfUeNgapId = 2

		UE = "eu"

	} else {
		mobileIdentity5GS = nasType.MobileIdentity5GS{
			Len:    12, // suci
			Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
		}
		ue.AuthenticationSubs = getAuthSubscription()
		ue.AmfUeNgapId = 1
		UE = "ue"
	}

	n3iwfUDPAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:500", ueContext.N3IWFIpAddress))
	if err != nil {
		pingLog.Fatal(err)
	}
	udpConnection := setupUDPSocket(ueContext, pingLog)

	initatorSPI, err := NewIKESecurityAssociation()
	if err != nil {
		pingLog.Fatal("Erro to generante InitiatorSPI ")
	} else {
		pingLog.Infof("New SPI generator %d \n\n", initatorSPI)
	}
	// IKE_SA_INIT
	ikeMessage := new(message.IKEMessage)
	ikeMessage.BuildIKEHeader(initatorSPI, 0, message.IKE_SA_INIT, message.InitiatorBitCheck, 0)

	// Security Association
	securityAssociation := ikeMessage.Payloads.BuildSecurityAssociation()
	// Proposal 1
	proposal := securityAssociation.Proposals.BuildProposal(1, message.TypeIKE, nil)
	// ENCR
	var attributeType uint16 = message.AttributeTypeKeyLength
	var keyLength uint16 = 256
	proposal.EncryptionAlgorithm.BuildTransform(message.TypeEncryptionAlgorithm, message.ENCR_AES_CBC, &attributeType, &keyLength, nil)
	// INTEG
	proposal.IntegrityAlgorithm.BuildTransform(message.TypeIntegrityAlgorithm, message.AUTH_HMAC_SHA1_96, nil, nil, nil)
	// PRF
	proposal.PseudorandomFunction.BuildTransform(message.TypePseudorandomFunction, message.PRF_HMAC_SHA1, nil, nil, nil)
	// DH
	proposal.DiffieHellmanGroup.BuildTransform(message.TypeDiffieHellmanGroup, message.DH_2048_BIT_MODP, nil, nil, nil)

	// Key exchange data
	generator := new(big.Int).SetUint64(handler.Group14Generator)
	factor, ok := new(big.Int).SetString(handler.Group14PrimeString, 16)
	if !ok {
		pingLog.Fatal("Generate key exchange datd failed")
	}
	secert := handler.GenerateRandomNumber()
	localPublicKeyExchangeValue := new(big.Int).Exp(generator, secert, factor).Bytes()
	prependZero := make([]byte, len(factor.Bytes())-len(localPublicKeyExchangeValue))
	localPublicKeyExchangeValue = append(prependZero, localPublicKeyExchangeValue...)
	ikeMessage.Payloads.BUildKeyExchange(message.DH_2048_BIT_MODP, localPublicKeyExchangeValue)

	// Nonce
	localNonce := handler.GenerateRandomNumber().Bytes()
	ikeMessage.Payloads.BuildNonce(localNonce)

	// Send to N3IWF
	ikeMessageData, err := ikeMessage.Encode()
	if err != nil {
		pingLog.Fatal(err)
	}
	if _, err := udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr); err != nil {
		pingLog.Fatal(err)
	}

	// Receive N3IWF reply
	buffer := make([]byte, 65535)
	n, _, err := udpConnection.ReadFromUDP(buffer)
	if err != nil {
		pingLog.Fatal(err)
	}
	ikeMessage.Payloads.Reset()
	err = ikeMessage.Decode(buffer[:n])
	if err != nil {
		pingLog.Fatal(err)
	}

	var sharedKeyExchangeData []byte
	var remoteNonce []byte

	for _, ikePayload := range ikeMessage.Payloads {
		switch ikePayload.Type() {
		case message.TypeSA:
			pingLog.Info("Get SA payload")
		case message.TypeKE:
			remotePublicKeyExchangeValue := ikePayload.(*message.KeyExchange).KeyExchangeData
			var i = 0
			for {
				if remotePublicKeyExchangeValue[i] != 0 {
					break
				}
			}
			remotePublicKeyExchangeValue = remotePublicKeyExchangeValue[i:]
			remotePublicKeyExchangeValueBig := new(big.Int).SetBytes(remotePublicKeyExchangeValue)
			sharedKeyExchangeData = new(big.Int).Exp(remotePublicKeyExchangeValueBig, secert, factor).Bytes()
		case message.TypeNiNr:
			remoteNonce = ikePayload.(*message.Nonce).NonceData
		}
	}

	ikeSecurityAssociation := &context.IKESecurityAssociation{
		LocalSPI:               initatorSPI,
		RemoteSPI:              ikeMessage.ResponderSPI,
		EncryptionAlgorithm:    proposal.EncryptionAlgorithm[0],
		IntegrityAlgorithm:     proposal.IntegrityAlgorithm[0],
		PseudorandomFunction:   proposal.PseudorandomFunction[0],
		DiffieHellmanGroup:     proposal.DiffieHellmanGroup[0],
		ConcatenatedNonce:      append(localNonce, remoteNonce...),
		DiffieHellmanSharedKey: sharedKeyExchangeData,
	}

	if err := generateKeyForIKESA(ikeSecurityAssociation); err != nil {
		pingLog.Fatalf("Generate key for IKE SA failed: %+v", err)
	}

	// IKE_AUTH
	ikeMessage.Payloads.Reset()
	ikeMessage.BuildIKEHeader(initatorSPI, ikeSecurityAssociation.RemoteSPI, message.IKE_AUTH, message.InitiatorBitCheck, 1)

	var ikePayload message.IKEPayloadContainer

	// Identification
	ikePayload.BuildIdentificationInitiator(message.ID_FQDN, []byte(UE))

	// Security Association
	securityAssociation = ikePayload.BuildSecurityAssociation()
	// Proposal 1
	proposal = securityAssociation.Proposals.BuildProposal(1, message.TypeESP, []byte{0, 0, 0, 1})
	// ENCR
	proposal.EncryptionAlgorithm.BuildTransform(message.TypeEncryptionAlgorithm, message.ENCR_AES_CBC, &attributeType, &keyLength, nil)
	// INTEG
	proposal.IntegrityAlgorithm.BuildTransform(message.TypeIntegrityAlgorithm, message.AUTH_HMAC_SHA1_96, nil, nil, nil)
	// ESN
	proposal.ExtendedSequenceNumbers.BuildTransform(message.TypeExtendedSequenceNumbers, message.ESN_NO, nil, nil, nil)

	// Traffic Selector
	tsi := ikePayload.BuildTrafficSelectorInitiator()
	tsi.TrafficSelectors.BuildIndividualTrafficSelector(message.TS_IPV4_ADDR_RANGE, 0, 0, 65535, []byte{0, 0, 0, 0}, []byte{255, 255, 255, 255})
	tsr := ikePayload.BuildTrafficSelectorResponder()
	tsr.TrafficSelectors.BuildIndividualTrafficSelector(message.TS_IPV4_ADDR_RANGE, 0, 0, 65535, []byte{0, 0, 0, 0}, []byte{255, 255, 255, 255})

	if err := encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage); err != nil {
		pingLog.Fatalf("Encrypting IKE message failed: %+v", err)
	}

	// Send to N3IWF
	ikeMessageData, err = ikeMessage.Encode()
	if err != nil {
		pingLog.Fatal(err)
	}
	if _, err := udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr); err != nil {
		pingLog.Fatal(err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		pingLog.Fatal(err)
	}
	ikeMessage.Payloads.Reset()
	err = ikeMessage.Decode(buffer[:n])
	if err != nil {
		pingLog.Fatal(err)
	}

	encryptedPayload, ok := ikeMessage.Payloads[0].(*message.Encrypted)
	if !ok {
		pingLog.Fatal("Received payload is not an encrypted payload")
	}

	decryptedIKEPayload, err := decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		pingLog.Fatalf("Decrypt IKE message failed: %+v", err)
	}

	var eapIdentifier uint8

	for _, ikePayload := range decryptedIKEPayload {
		switch ikePayload.Type() {
		case message.TypeIDr:
			pingLog.Info("Get IDr")
		case message.TypeAUTH:
			pingLog.Info("Get AUTH")
		case message.TypeCERT:
			pingLog.Info("Get CERT")
		case message.TypeEAP:
			eapIdentifier = ikePayload.(*message.EAP).Identifier
			pingLog.Info("Get EAP")
		}
	}

	// IKE_AUTH - EAP exchange
	ikeMessage.Payloads.Reset()
	ikeMessage.BuildIKEHeader(initatorSPI, ikeSecurityAssociation.RemoteSPI, message.IKE_AUTH, message.InitiatorBitCheck, 2)

	ikePayload.Reset()

	// EAP-5G vendor type data
	eapVendorTypeData := make([]byte, 2)
	eapVendorTypeData[0] = message.EAP5GType5GNAS

	// AN Parameters
	anParameters := buildEAP5GANParameters()
	anParametersLength := make([]byte, 2)
	binary.BigEndian.PutUint16(anParametersLength, uint16(len(anParameters)))
	eapVendorTypeData = append(eapVendorTypeData, anParametersLength...)
	eapVendorTypeData = append(eapVendorTypeData, anParameters...)

	// NAS
	registrationRequest := ue_nas.RegistrationRequest(
		nasMessage.RegistrationType5GSInitialRegistration,
		mobileIdentity5GS,
		nil,
		ue.GetUESecurityCapability(),
		nil,
		nil,
		nil,
	)

	nasLength := make([]byte, 2)
	binary.BigEndian.PutUint16(nasLength, uint16(len(registrationRequest)))
	eapVendorTypeData = append(eapVendorTypeData, nasLength...)
	eapVendorTypeData = append(eapVendorTypeData, registrationRequest...)

	eap := ikePayload.BuildEAP(message.EAPCodeResponse, eapIdentifier)
	eap.EAPTypeData.BuildEAPExpanded(message.VendorID3GPP, message.VendorTypeEAP5G, eapVendorTypeData)

	if err := encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage); err != nil {
		pingLog.Fatal(err)
	}

	// Send to N3IWF
	ikeMessageData, err = ikeMessage.Encode()
	if err != nil {
		pingLog.Fatal(err)
	}
	if _, err := udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr); err != nil {
		pingLog.Fatal(err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		pingLog.Fatal(err)
	}
	ikeMessage.Payloads.Reset()
	err = ikeMessage.Decode(buffer[:n])
	if err != nil {
		pingLog.Fatal(err)
	}
	encryptedPayload, ok = ikeMessage.Payloads[0].(*message.Encrypted)
	if !ok {
		pingLog.Fatal("Received payload is not an encrypted payload")
	}
	decryptedIKEPayload, err = decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		pingLog.Fatalf("Decrypt IKE message failed: %+v", err)
	}

	var eapReq *message.EAP
	var eapExpanded *message.EAPExpanded

	eapReq, ok = decryptedIKEPayload[0].(*message.EAP)
	if !ok {
		pingLog.Fatal("Received packet is not an EAP payload")
	}

	var decodedNAS *nas.Message

	eapExpanded, ok = eapReq.EAPTypeData[0].(*message.EAPExpanded)
	if !ok {
		pingLog.Fatal("The EAP data is not an EAP expended.")
	}

	// Decode NAS - Authentication Request
	nasData := eapExpanded.VendorData[4:]
	decodedNAS = new(nas.Message)
	if err := decodedNAS.PlainNasDecode(&nasData); err != nil {
		pingLog.Fatal(err)
	}

	// Calculate for RES*
	rand := decodedNAS.AuthenticationRequest.GetRANDValue()
	resStat := ue.DeriveRESstarAndSetKey(ue.AuthenticationSubs, rand[:], "5G:mnc093.mcc208.3gppnetwork.org")

	// send NAS Authentication Response
	pdu := ue_nas.AuthenticationResponse(resStat, "")

	// IKE_AUTH - EAP exchange
	ikeMessage.Payloads.Reset()
	ikeMessage.BuildIKEHeader(initatorSPI, ikeSecurityAssociation.RemoteSPI, message.IKE_AUTH, message.InitiatorBitCheck, 3)

	ikePayload.Reset()

	// EAP-5G vendor type data
	eapVendorTypeData = make([]byte, 4)
	eapVendorTypeData[0] = message.EAP5GType5GNAS

	// NAS - Authentication Response
	nasLength = make([]byte, 2)
	binary.BigEndian.PutUint16(nasLength, uint16(len(pdu)))
	eapVendorTypeData = append(eapVendorTypeData, nasLength...)
	eapVendorTypeData = append(eapVendorTypeData, pdu...)

	eap = ikePayload.BuildEAP(message.EAPCodeResponse, eapReq.Identifier)
	eap.EAPTypeData.BuildEAPExpanded(message.VendorID3GPP, message.VendorTypeEAP5G, eapVendorTypeData)

	err = encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage)
	if err != nil {
		pingLog.Fatal(err)
	}

	// Send to N3IWF
	ikeMessageData, err = ikeMessage.Encode()
	if err != nil {
		pingLog.Fatal(err)
	}
	_, err = udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr)
	if err != nil {
		pingLog.Fatal(err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		pingLog.Fatal(err)
	}
	ikeMessage.Payloads.Reset()
	err = ikeMessage.Decode(buffer[:n])
	if err != nil {
		pingLog.Fatal(err)
	}
	encryptedPayload, ok = ikeMessage.Payloads[0].(*message.Encrypted)
	if !ok {
		pingLog.Fatal("Received pakcet is not and encrypted payload")
	}
	decryptedIKEPayload, err = decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		pingLog.Fatal(err)
	}
	eapReq, ok = decryptedIKEPayload[0].(*message.EAP)
	if !ok {
		pingLog.Fatal("Received packet is not an EAP payload")
	}
	eapExpanded, ok = eapReq.EAPTypeData[0].(*message.EAPExpanded)
	if !ok {
		pingLog.Fatal("Received packet is not an EAP expended payload")
	}

	nasData = eapExpanded.VendorData[4:]

	// Send NAS Security Mode Complete Msg
	registrationRequestWith5GMM := ue_nas.RegistrationRequest(
		nasMessage.RegistrationType5GSInitialRegistration,
		mobileIdentity5GS,
		nil,
		ue.GetUESecurityCapability(),
		nil,
		nil,
		nil,
	)
	pdu = ue_nas.SecurityModeComplete(registrationRequestWith5GMM)
	pdu, err = EncodeNasPduWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext, true, true)
	//assert.Nil(t, err)
	if err != nil {
		pingLog.Error("Erro na função que que code pdu session ")
	}

	// IKE_AUTH - EAP exchange
	ikeMessage.Payloads.Reset()
	ikeMessage.BuildIKEHeader(initatorSPI, ikeSecurityAssociation.RemoteSPI, message.IKE_AUTH, message.InitiatorBitCheck, 4)

	ikePayload.Reset()

	// EAP-5G vendor type data
	eapVendorTypeData = make([]byte, 4)
	eapVendorTypeData[0] = message.EAP5GType5GNAS

	// NAS - Authentication Response
	nasLength = make([]byte, 2)
	binary.BigEndian.PutUint16(nasLength, uint16(len(pdu)))
	eapVendorTypeData = append(eapVendorTypeData, nasLength...)
	eapVendorTypeData = append(eapVendorTypeData, pdu...)

	eap = ikePayload.BuildEAP(message.EAPCodeResponse, eapReq.Identifier)
	eap.EAPTypeData.BuildEAPExpanded(message.VendorID3GPP, message.VendorTypeEAP5G, eapVendorTypeData)

	err = encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage)
	if err != nil {
		pingLog.Fatal(err)
	}

	// Send to N3IWF
	ikeMessageData, err = ikeMessage.Encode()
	if err != nil {
		pingLog.Fatal(err)
	}
	_, err = udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr)
	if err != nil {
		pingLog.Fatal(err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		pingLog.Fatal(err)
	}
	ikeMessage.Payloads.Reset()
	err = ikeMessage.Decode(buffer[:n])
	if err != nil {
		pingLog.Fatal(err)
	}
	encryptedPayload, ok = ikeMessage.Payloads[0].(*message.Encrypted)
	if !ok {
		pingLog.Fatal("Received pakcet is not and encrypted payload")
	}
	decryptedIKEPayload, err = decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		pingLog.Fatal(err)
	}
	eapReq, ok = decryptedIKEPayload[0].(*message.EAP)
	if !ok {
		pingLog.Fatal("Received packet is not an EAP payload")
	}
	if eapReq.Code != message.EAPCodeSuccess {
		pingLog.Fatal("Not Success")
	}

	// IKE_AUTH - Authentication
	ikeMessage.Payloads.Reset()
	ikeMessage.BuildIKEHeader(initatorSPI, ikeSecurityAssociation.RemoteSPI, message.IKE_AUTH, message.InitiatorBitCheck, 5)

	ikePayload.Reset()

	// Authentication
	ikePayload.BuildAuthentication(message.SharedKeyMesageIntegrityCode, []byte{1, 2, 3})

	// Configuration Request
	configurationRequest := ikePayload.BuildConfiguration(message.CFG_REQUEST)
	configurationRequest.ConfigurationAttribute.BuildConfigurationAttribute(message.INTERNAL_IP4_ADDRESS, nil)

	err = encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage)
	if err != nil {
		pingLog.Fatal(err)
	}

	// Send to N3IWF
	ikeMessageData, err = ikeMessage.Encode()
	if err != nil {
		pingLog.Fatal(err)
	}
	_, err = udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr)
	if err != nil {
		pingLog.Fatal(err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		pingLog.Fatal(err)
	}
	ikeMessage.Payloads.Reset()
	err = ikeMessage.Decode(buffer[:n])
	if err != nil {
		pingLog.Fatal(err)
	}
	encryptedPayload, ok = ikeMessage.Payloads[0].(*message.Encrypted)
	if !ok {
		pingLog.Fatal("Received pakcet is not and encrypted payload")
	}
	decryptedIKEPayload, err = decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		pingLog.Fatal(err)
	}

	// AUTH, SAr2, TSi, Tsr, N(NAS_IP_ADDRESS), N(NAS_TCP_PORT)
	var responseSecurityAssociation *message.SecurityAssociation
	var responseTrafficSelectorInitiator *message.TrafficSelectorInitiator
	var responseTrafficSelectorResponder *message.TrafficSelectorResponder
	var responseConfiguration *message.Configuration
	n3iwfNASAddr := new(net.TCPAddr)
	ueAddr := new(net.IPNet)

	for _, ikePayload := range decryptedIKEPayload {
		switch ikePayload.Type() {
		case message.TypeAUTH:
			pingLog.Info("Get Authentication from N3IWF")
		case message.TypeSA:
			responseSecurityAssociation = ikePayload.(*message.SecurityAssociation)
			ikeSecurityAssociation.IKEAuthResponseSA = responseSecurityAssociation
		case message.TypeTSi:
			responseTrafficSelectorInitiator = ikePayload.(*message.TrafficSelectorInitiator)
		case message.TypeTSr:
			responseTrafficSelectorResponder = ikePayload.(*message.TrafficSelectorResponder)
		case message.TypeN:
			notification := ikePayload.(*message.Notification)
			if notification.NotifyMessageType == message.Vendor3GPPNotifyTypeNAS_IP4_ADDRESS {
				n3iwfNASAddr.IP = net.IPv4(notification.NotificationData[0], notification.NotificationData[1], notification.NotificationData[2], notification.NotificationData[3])
			}
			if notification.NotifyMessageType == message.Vendor3GPPNotifyTypeNAS_TCP_PORT {
				n3iwfNASAddr.Port = int(binary.BigEndian.Uint16(notification.NotificationData))
			}
		case message.TypeCP:
			responseConfiguration = ikePayload.(*message.Configuration)
			if responseConfiguration.ConfigurationType == message.CFG_REPLY {
				for _, configAttr := range responseConfiguration.ConfigurationAttribute {
					if configAttr.Type == message.INTERNAL_IP4_ADDRESS {
						ueAddr.IP = configAttr.Value
					}
					if configAttr.Type == message.INTERNAL_IP4_NETMASK {
						ueAddr.Mask = configAttr.Value
					}
				}
			}
		}
	}

	childSecurityAssociationContext, err := createIKEChildSecurityAssociation(ikeSecurityAssociation.IKEAuthResponseSA)
	if err != nil {
		pingLog.Fatalf("Create child security association context failed: %+v", err)
		return
	}
	err = parseIPAddressInformationToChildSecurityAssociation(childSecurityAssociationContext, net.ParseIP(ueContext.N3IWFIpAddress), ueContext.IKEBindAddress, responseTrafficSelectorInitiator.TrafficSelectors[0], responseTrafficSelectorResponder.TrafficSelectors[0])
	if err != nil {
		pingLog.Fatalf("Parse IP address to child security association failed: %+v", err)
		return
	}
	// Select TCP traffic
	childSecurityAssociationContext.SelectedIPProtocol = unix.IPPROTO_TCP

	if err := generateKeyForChildSA(ikeSecurityAssociation, childSecurityAssociationContext); err != nil {
		pingLog.Fatalf("Generate key for child SA failed: %+v", err)
		return
	}

	// Aplly XFRM rules
	if err = applyXFRMRule(true, childSecurityAssociationContext); err != nil {
		pingLog.Fatalf("Applying XFRM rules failed: %+v", err)
		return
	}

	// Get link ipsec0
	links, err := netlink.LinkList()
	if err != nil {
		pingLog.Fatal(err)
	}

	var linkIPSec netlink.Link
	for _, link := range links {
		if link.Attrs() != nil {
			if link.Attrs().Name == "ipsec0" {
				linkIPSec = link
				break
			}
		}
	}
	if linkIPSec == nil {
		pingLog.Fatal("No link named ipsec0")
	}

	linkIPSecAddr := &netlink.Addr{
		IPNet: ueAddr,
	}

	if err := netlink.AddrAdd(linkIPSec, linkIPSecAddr); err != nil {
		pingLog.Fatalf("Set ipsec0 addr failed: %v", err)
	}

	defer func() {
		_ = netlink.AddrDel(linkIPSec, linkIPSecAddr)
		// NOTE (Matan): The XfrmPolicyFlush() function flags on Mac
		// My theory is that since it is sourced from netlink's xfrm_policy_linux, it may not be exposed on Mac
		_ = netlink.XfrmPolicyFlush()
		_ = netlink.XfrmStateFlush(netlink.XFRM_PROTO_IPSEC_ANY)
	}()

	localTCPAddr := &net.TCPAddr{
		IP: ueAddr.IP,
	}
	tcpConnWithN3IWF, err := net.DialTCP("tcp", localTCPAddr, n3iwfNASAddr)
	if err != nil {
		pingLog.Fatal(err)
	}

	nasMsg := make([]byte, 65535)

	_, err = tcpConnWithN3IWF.Read(nasMsg)
	if err != nil {
		pingLog.Fatal(err)
	}

	// send NAS Registration Complete Msg
	pdu = ue_nas.RegistrationComplete(nil)
	pdu, err = EncodeNasPduWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true, false)
	if err != nil {
		pingLog.Fatal(err)
	}
	_, err = tcpConnWithN3IWF.Write(pdu)
	if err != nil {
		pingLog.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	// UE request PDU session setup
	sNssai := models.Snssai{
		Sst: 1,
		Sd:  "010203",
	}
	// TODO: why is this placed in the initial registration? Should be in the PDU session initializing func
	pdu = ue_nas.ULNASTransportPDUSessionEstablishmentRequest(ueContext.PDUSessionID, nasMessage.ULNASTransportRequestTypeInitialRequest, ueContext.Dnn, &sNssai)
	pdu, err = EncodeNasPduWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true, false)
	if err != nil {
		pingLog.Fatal(err)
	}
	_, err = tcpConnWithN3IWF.Write(pdu)
	if err != nil {
		pingLog.Fatal(err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		pingLog.Fatal(err)
	}
	ikeMessage.Payloads.Reset()
	err = ikeMessage.Decode(buffer[:n])
	if err != nil {
		pingLog.Fatal(err)
	}

	pingLog.Infof("IKE message exchange type: %d", ikeMessage.ExchangeType)
	pingLog.Infof("IKE message ID: %d", ikeMessage.MessageID)
	encryptedPayload, ok = ikeMessage.Payloads[0].(*message.Encrypted)
	if !ok {
		pingLog.Fatal("Received pakcet is not and encrypted payload")
	}
	decryptedIKEPayload, err = decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		pingLog.Fatal(err)
	}

	var upIPAddr net.IP
	for _, ikePayload := range decryptedIKEPayload {
		switch ikePayload.Type() {
		case message.TypeSA:
			responseSecurityAssociation = ikePayload.(*message.SecurityAssociation)
		case message.TypeTSi:
			responseTrafficSelectorInitiator = ikePayload.(*message.TrafficSelectorInitiator)
		case message.TypeTSr:
			responseTrafficSelectorResponder = ikePayload.(*message.TrafficSelectorResponder)
		case message.TypeN:
			notification := ikePayload.(*message.Notification)
			if notification.NotifyMessageType == message.Vendor3GPPNotifyType5G_QOS_INFO {
				pingLog.Info("Received Qos Flow settings")
			}
			if notification.NotifyMessageType == message.Vendor3GPPNotifyTypeUP_IP4_ADDRESS {
				pingLog.Info("UP IP Address: %+v\n", notification.NotificationData)
				upIPAddr = notification.NotificationData[:4]
			}
		case message.TypeNiNr:
			responseNonce := ikePayload.(*message.Nonce)
			ikeSecurityAssociation.ConcatenatedNonce = responseNonce.NonceData
		}
	}

	// IKE CREATE_CHILD_SA response
	ikeMessage.Payloads.Reset()
	ikeMessage.BuildIKEHeader(ikeMessage.InitiatorSPI, ikeMessage.ResponderSPI, message.CREATE_CHILD_SA, message.ResponseBitCheck, ikeMessage.MessageID)

	ikePayload.Reset()

	// SA
	ikePayload = append(ikePayload, responseSecurityAssociation)

	// TSi
	ikePayload = append(ikePayload, responseTrafficSelectorInitiator)

	// TSr
	ikePayload = append(ikePayload, responseTrafficSelectorResponder)

	// Nonce
	localNonce = handler.GenerateRandomNumber().Bytes()
	ikeSecurityAssociation.ConcatenatedNonce = append(ikeSecurityAssociation.ConcatenatedNonce, localNonce...)
	ikePayload.BuildNonce(localNonce)

	if err := encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage); err != nil {
		pingLog.Fatal(err)
	}

	// Send to N3IWF
	ikeMessageData, err = ikeMessage.Encode()
	if err != nil {
		pingLog.Fatal(err)
	}
	_, err = udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr)
	if err != nil {
		pingLog.Fatal(err)
	}

	childSecurityAssociationContextUserPlane, err := createIKEChildSecurityAssociation(responseSecurityAssociation)
	if err != nil {
		pingLog.Fatalf("Create child security association context failed: %+v", err)
		return
	}
	err = parseIPAddressInformationToChildSecurityAssociation(childSecurityAssociationContextUserPlane, net.ParseIP(ueContext.N3IWFIpAddress), ueContext.IKEBindAddress, responseTrafficSelectorResponder.TrafficSelectors[0], responseTrafficSelectorInitiator.TrafficSelectors[0])
	if err != nil {
		pingLog.Fatalf("Parse IP address to child security association failed: %+v", err)
		return
	}
	// Select GRE traffic
	childSecurityAssociationContextUserPlane.SelectedIPProtocol = unix.IPPROTO_GRE

	if err := generateKeyForChildSA(ikeSecurityAssociation, childSecurityAssociationContextUserPlane); err != nil {
		pingLog.Fatalf("Generate key for child SA failed: %+v", err)
		return
	}

	pingLog.Infof("State function: encr: %d, auth: %d", childSecurityAssociationContextUserPlane.EncryptionAlgorithm, childSecurityAssociationContextUserPlane.IntegrityAlgorithm)
	// Aplly XFRM rules
	if err = applyXFRMRule(false, childSecurityAssociationContextUserPlane); err != nil {
		pingLog.Fatalf("Applying XFRM rules failed: %+v", err)
		return
	}

	// New GRE tunnel interface
	newGRETunnel := &netlink.Gretun{
		LinkAttrs: netlink.LinkAttrs{
			Name: "gretun0",
		},
		Local:  ueAddr.IP,
		Remote: upIPAddr,
	}
	if err := netlink.LinkAdd(newGRETunnel); err != nil {
		pingLog.Fatal(err)
	}
	// Get link info
	links, err = netlink.LinkList()
	if err != nil {
		pingLog.Fatal(err)
	}
	var linkGRE netlink.Link
	for _, link := range links {
		if link.Attrs() != nil {
			if link.Attrs().Name == "gretun0" {
				linkGRE = link
				break
			}
		}
	}
	if linkGRE == nil {
		pingLog.Fatal("No link named gretun0")
	}

	// Link address 60.60.0.1/24
	ipv4 := net.ParseIP(ueContext.GREIPAddress)
	linkGREAddr := &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   ipv4,
			Mask: net.IPv4Mask(255, 255, 255, 255),
		},
	}
	if err := netlink.AddrAdd(linkGRE, linkGREAddr); err != nil {
		pingLog.Fatal(err)
	}
	// Set GRE interface up
	if err := netlink.LinkSetUp(linkGRE); err != nil {
		pingLog.Fatal(err)
	}
	// Add route
	upRoute := &netlink.Route{
		LinkIndex: linkGRE.Attrs().Index,
		Dst: &net.IPNet{
			IP:   net.IPv4zero,
			Mask: net.IPv4Mask(0, 0, 0, 0),
		},
	}
	if err := netlink.RouteAdd(upRoute); err != nil {
		pingLog.Fatal(err)
	}

	defer func() {
		_ = netlink.LinkSetDown(linkGRE)
		_ = netlink.LinkDel(linkGRE)
	}()

	time.Sleep(300 * time.Hour)

}

func setUESecurityCapability(ue *RanUeContext) (UESecurityCapability *nasType.UESecurityCapability) {
	UESecurityCapability = &nasType.UESecurityCapability{
		Iei:    nasMessage.RegistrationRequestUESecurityCapabilityType,
		Len:    8,
		Buffer: []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	switch ue.CipheringAlg {
	case security.AlgCiphering128NEA0:
		UESecurityCapability.SetEA0_5G(1)
	case security.AlgCiphering128NEA1:
		UESecurityCapability.SetEA1_128_5G(1)
	case security.AlgCiphering128NEA2:
		UESecurityCapability.SetEA2_128_5G(1)
	case security.AlgCiphering128NEA3:
		UESecurityCapability.SetEA3_128_5G(1)
	}

	switch ue.IntegrityAlg {
	case security.AlgIntegrity128NIA0:
		UESecurityCapability.SetIA0_5G(1)
	case security.AlgIntegrity128NIA1:
		UESecurityCapability.SetIA1_128_5G(1)
	case security.AlgIntegrity128NIA2:
		UESecurityCapability.SetIA2_128_5G(1)
	case security.AlgIntegrity128NIA3:
		UESecurityCapability.SetIA3_128_5G(1)
	}

	return
}

func HandleRegistrationProcedure(ueContext *ue_context.UEContext) {

	// TODO: #LABORA what is the current step??
	/*
		if IDLE and DEREGISTERED
			execute Initial Registration Procedure
		elif IDLE and REGISTERED
			execute Service Request
		elif CONNECTED AND DEREGISTERED
			execute Registration Procedure
		else
			execute Registration Procedure
	*/
	InitialRegistrationProcedure(ueContext)
}

func HandleDeregistrationProcedure(ueContext *ue_context.UEContext) {

}

func SetupPDUSession(ueContext *ue_context.UEContext) {

}

func DeregistrationProcedure(ueContext *ue_context.UEContext) {

}
