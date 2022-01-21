package ue_util

import (
	"github.com/free5gc/path_util"
)

var UeLogPath = path_util.Free5gcPath("free5gc/uesslkey.log")

// UePemPath TODO: #LABORA Create the key and pem certificate files for UE in free5gc/support
var UePemPath = path_util.Free5gcPath("free5gc/support/TLS/ue.pem")
var UeKeyPath = path_util.Free5gcPath("free5gc/support/TLS/ue.key")

// DefaultUeConfigPath TODO: #LABORA Add free5gc/config/uecfg.conf
var DefaultUeConfigPath = path_util.Free5gcPath("free5gc/config/uecfg.conf")
