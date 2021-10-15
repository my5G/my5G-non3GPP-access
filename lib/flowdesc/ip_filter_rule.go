package flowdesc

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

const ProtocalNumberAny = 0xfc

// Action - Action of IPFilterRule
type Action string

// Action const
const (
	Permit Action = "permit"
	Deny   Action = "deny"
)

// Direction - direction of IPFilterRule
type Direction string

// Direction const
const (
	In  Direction = "in"
	Out Direction = "out"
)

func flowDescErrorf(format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	return fmt.Errorf("flowdesc: %s", msg)
}

// IPFilterRule define RFC 3588 that referd by TS 29.212
type IPFilterRule struct {
	action   Action
	dir      Direction
	proto    uint8  // protocal number
	srcIP    string // <address/mask>
	srcPorts string // [ports]
	dstIP    string // <address/mask>
	dstPorts string // [ports]
}

// NewIPFilterRule returns a new IPFilterRule instance
func NewIPFilterRule() *IPFilterRule {
	r := &IPFilterRule{
		action:   Permit,
		dir:      Out,
		proto:    ProtocalNumberAny,
		srcIP:    "",
		srcPorts: "",
		dstIP:    "",
		dstPorts: "",
	}
	return r
}

// SetAction sets action of the IPFilterRule
func (r *IPFilterRule) SetAction(action Action) error {
	switch action {
	case Permit:
		r.action = action
	case Deny:
		r.action = action
	default:
		return flowDescErrorf("'%s' is not allow, action only accept 'permit' or 'deny'", action)
	}
	return nil
}

// GetAction returns action of the IPFilterRule
func (r *IPFilterRule) GetAction() Action {
	return r.action
}

// SetDirection sets direction of the IPFilterRule
func (r *IPFilterRule) SetDirection(dir Direction) error {
	switch dir {
	case Out:
		r.dir = dir
	case In:
		return flowDescErrorf("dir cannot be 'in' in core-network")
	default:
		return flowDescErrorf("'%s' is not allow, dir only accept 'out'", dir)
	}
	return nil
}

// GetDirection returns direction of the IPFilterRule
func (r *IPFilterRule) GetDirection() Direction {
	return r.dir
}

// SetProtocal sets IP protocal number of the IPFilterRule
// 0xfc stand for ip (any)
func (r *IPFilterRule) SetProtocal(proto uint8) error {
	r.proto = proto
	return nil
}

// GetProtocal returns the ip protocal number of the IPFilterRule
func (r *IPFilterRule) GetProtocal() uint8 {
	return r.proto
}

// SetSourceIP sets source IP of the IPFilterRule
func (r *IPFilterRule) SetSourceIP(networkStr string) error {
	if networkStr == "any" || networkStr == "assigned" {
		r.srcIP = networkStr
		return nil
	}
	if networkStr[0] == '!' {
		return flowDescErrorf("Base on TS 29.212, ! expression shall not be used")
	}

	var ipStr string

	ip := net.ParseIP(networkStr)
	if ip == nil {
		_, ipNet, err := net.ParseCIDR(networkStr)
		if err != nil {
			return flowDescErrorf("Source IP format error")
		}
		ipStr = ipNet.String()
	} else {
		ipStr = ip.String()
	}

	r.srcIP = ipStr
	return nil
}

// GetSourceIP returns src of the IPFilterRule
func (r *IPFilterRule) GetSourceIP() string {
	return r.srcIP
}

// SetSourcePorts sets source ports of the IPFilterRule
func (r *IPFilterRule) SetSourcePorts(ports string) error {

	if ports == "" {
		r.srcPorts = ""
		return nil
	}

	if match, err := regexp.MatchString("^[0-9]+(-[0-9]+)?(,[0-9]+)*$", ports); err != nil || !match {
		return flowDescErrorf("not valid format of port number")
	}

	// Check port range
	portSlice := regexp.MustCompile(`[\\,\\-]+`).Split(ports, -1)
	for _, portStr := range portSlice {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return err
		}
		if 0 > int(port) || port > 65535 {
			return errors.New("Invalid port number")
		}
	}

	r.srcPorts = ports
	return nil
}

// GetSourcePorts returns src ports of the IPFilterRule
func (r *IPFilterRule) GetSourcePorts() string {
	return r.srcPorts
}

// SetDestinationIP sets destination IP of the IPFilterRule
func (r *IPFilterRule) SetDestinationIP(networkStr string) error {
	if networkStr == "any" || networkStr == "assigned" {
		r.dstIP = networkStr
		return nil
	}
	if networkStr[0] == '!' {
		return flowDescErrorf("Base on TS 29.212, ! expression shall not be used")
	}

	var ipDst string

	ip := net.ParseIP(networkStr)
	if ip == nil {
		_, ipNet, err := net.ParseCIDR(networkStr)
		if err != nil {
			return flowDescErrorf("Source IP format error")
		}
		ipDst = ipNet.String()
	} else {
		ipDst = ip.String()
	}

	r.dstIP = ipDst
	return nil
}

// GetDestinationIP returns dst of the IPFilterRule
func (r *IPFilterRule) GetDestinationIP() string {
	return r.dstIP
}

// SetDestinationPorts sets destination ports of the IPFilterRule
func (r *IPFilterRule) SetDestinationPorts(ports string) error {

	if ports == "" {
		r.dstPorts = ports
		return nil
	}

	match, err := regexp.MatchString("^[0-9]+(-[0-9]+)?(,[0-9]+)*$", ports)

	if err != nil {
		return flowDescErrorf("Regex match error")
	}
	if !match {
		return flowDescErrorf("Ports format error")
	}

	// Check port range
	portSlice := regexp.MustCompile(`[\\,\\-]+`).Split(ports, -1)
	for _, portStr := range portSlice {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return err
		}
		if 0 > int(port) || port > 65535 {
			return flowDescErrorf("Invalid port number")
		}
	}

	r.dstPorts = ports
	return nil
}

// GetDestinationPorts returns src ports of the IPFilterRule
func (r *IPFilterRule) GetDestinationPorts() string {
	return r.dstPorts
}

// SwapSourceAndDestination swap the src and dst of the IPFilterRule
func (r *IPFilterRule) SwapSourceAndDestination() {
	r.srcIP, r.dstIP = r.dstIP, r.srcIP
	r.srcPorts, r.dstPorts = r.dstPorts, r.srcPorts
}

// Encode function out put the IPFilterRule from the struct
func Encode(r *IPFilterRule) (string, error) {
	var ipFilterRuleStr []string

	// pre-allocate seven element
	ipFilterRuleStr = make([]string, 0, 9)

	// action
	switch r.action {
	case Permit:
		ipFilterRuleStr = append(ipFilterRuleStr, "permit")
	case Deny:
		ipFilterRuleStr = append(ipFilterRuleStr, "deny")
	}

	// dir
	switch r.dir {
	case Out:
		ipFilterRuleStr = append(ipFilterRuleStr, "out")
	}

	// proto
	if r.proto == ProtocalNumberAny {
		ipFilterRuleStr = append(ipFilterRuleStr, "ip")
	} else {
		ipFilterRuleStr = append(ipFilterRuleStr, strconv.Itoa(int(r.proto)))
	}

	// from
	ipFilterRuleStr = append(ipFilterRuleStr, "from")

	// src
	if r.srcIP != "" {
		ipFilterRuleStr = append(ipFilterRuleStr, r.srcIP)
	} else {
		return "", errors.New("Encode error: no source IP")
	}
	if r.srcPorts != "" {
		ipFilterRuleStr = append(ipFilterRuleStr, r.srcPorts)
	}

	// to
	ipFilterRuleStr = append(ipFilterRuleStr, "to")

	// dst
	if r.dstIP != "" {
		ipFilterRuleStr = append(ipFilterRuleStr, r.dstIP)
	} else {
		return "", errors.New("Encode error: no destination IP")
	}
	if r.dstPorts != "" {
		ipFilterRuleStr = append(ipFilterRuleStr, r.dstPorts)
	}

	// accroding TS 29.212 IPFilterRule cannot use [options]

	return strings.Join(ipFilterRuleStr, " "), nil
}

// Decode parsing the string to IPFilterRule
func Decode(s string, r *IPFilterRule) error {
	parts := strings.Split(s, " ")

	var ptr int

	// action
	if err := r.SetAction(Action(parts[ptr])); err != nil {
		return err
	}
	ptr++

	// dir
	if err := r.SetDirection(Direction(parts[ptr])); err != nil {
		return err
	}
	ptr++

	// proto
	var protoNumber uint8
	if parts[ptr] == "ip" {
		r.proto = ProtocalNumberAny
	} else {
		if proto, err := strconv.Atoi(parts[ptr]); err != nil {
			return flowDescErrorf("parse proto failed: %s", err)
		} else {
			protoNumber = uint8(proto)
		}
		if err := r.SetProtocal(protoNumber); err != nil {
			return flowDescErrorf("parse proto failed: %s", err)
		}
	}
	ptr++

	// from
	if from := parts[ptr]; from != "from" {
		return flowDescErrorf("parse faild: must have 'from'")
	}
	ptr++

	// src
	if err := r.SetSourceIP(parts[ptr]); err != nil {
		return err
	}
	ptr++

	if err := r.SetSourcePorts(parts[ptr]); err != nil {
	} else {
		ptr++
	}

	// to
	if to := parts[ptr]; to != "to" {
		return flowDescErrorf("parse faild: must have 'to'")
	}
	ptr++

	// dst
	if err := r.SetDestinationIP(parts[ptr]); err != nil {
		return err
	}
	ptr++

	// if end of parts
	if !(len(parts) > ptr) {
		return nil
	}

	if err := r.SetDestinationPorts(parts[ptr]); err != nil {
		return err
	} //else {
	//ptr++
	//}

	return nil
}
